package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mail/internal/config"
	"mail/internal/domain/models"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Client struct {
	Reader                           *kafka.Reader
	logger                           *zap.SugaredLogger
	workersCount                     int
	eventsToProcessChannelBufferSize int
	processedEventsChannelBufferSize int
	mailRateLimitSec                 int
}

type mail struct {
	kafkaMessage kafka.Message
	mailResult   string
}

func New(config *config.Config, logger *zap.SugaredLogger) *Client {
	c := Client{}
	brokers := []string{config.KafkaHost + ":" + config.KafkaPort}
	topic := config.KafkaTopic
	groupID := "mail_group"

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})
	c.logger = logger
	c.workersCount = config.WorkersCount
	c.eventsToProcessChannelBufferSize = config.EventsToProcessChannelBufferSize
	c.processedEventsChannelBufferSize = config.ProcessedEventsChannelBufferSize
	c.mailRateLimitSec = config.MailRateLimitSec

	return &c
}

func (c *Client) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}

	eventsToProcess := make(chan kafka.Message, c.eventsToProcessChannelBufferSize)
	processedEvents := make(chan mail, c.processedEventsChannelBufferSize)

	for i := 0; i < c.workersCount; i++ {
		wg.Add(1)
		c.logger.Infof("started worker: %d", i)

		go func() {
			defer wg.Done()
			c.worker(ctx, eventsToProcess, processedEvents)
		}()
	}

	go func() {
		for {
			msg, err := c.Reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.logger.Infof("kafka reader stopped by context: %v", err)

					return
				}

				c.logger.Errorf("kafka reader failed: %v", err)

				continue
			}

			c.logger.Infof("fetched message: %v value: %s, ofset: %v", msg.Key, msg.Value, msg.Offset)
			eventsToProcess <- msg
		}
	}()

	go func() {
		wg.Wait()
		close(processedEvents)
	}()

	for res := range processedEvents {
		// Here can be work with SMTP server
		c.logger.Infof(res.mailResult)

		err := c.Reader.CommitMessages(ctx, res.kafkaMessage)
		if err != nil {
			c.logger.Errorf("kafka commit message failed: %v. Message: %v", err, res.kafkaMessage)
		}

		// Implementation of RATE LIMIT
		time.Sleep(time.Duration(c.mailRateLimitSec) * time.Second)
	}
}

func (c *Client) Stop() {
	err := c.Reader.Close()
	if err != nil {
		c.logger.Errorf("kafka reader close failed: %v", err)
	}

	c.logger.Info("kafka reader closed")
}

func (c *Client) worker(ctx context.Context, toProcess <-chan kafka.Message, processed chan<- mail) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("worker was finished")

			return
		case value, ok := <-toProcess:
			if !ok {
				return
			}

			var mailEvent models.MailEvent

			err := json.Unmarshal(value.Value, &mailEvent)
			if err != nil {
				c.logger.Errorf("kafka message unmarshall failed: %v, value: %s", err, value.Value)
			}

			// Assume a worker needs 10 sec to process mailEvent
			time.Sleep(10 * time.Second)

			mailRes := fmt.Sprintf(
				`
						Subject: %s
						Addresse: %s
						Body: 
						Please approve following.
						%s
						YES: %s
						NO: %s`,
				mailEvent.Description, mailEvent.Addressee, mailEvent.Body,
				mailEvent.ApproveLink, mailEvent.RejectLink)

			processed <- mail{kafkaMessage: value, mailResult: mailRes}

			c.logger.Infof("worker: email created for message %d", value.Offset)
		}
	}
}
