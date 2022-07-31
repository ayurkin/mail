package application

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mail/internal/adapters/kafka"
	"mail/internal/config"
)

var (
	logger      *zap.Logger
	mailHandler *consumer.Client
)

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()

	appConfig, err := config.NewConfig()
	if err != nil {
		logger.Sugar().Fatalf("create config failed: %v", err)
	}

	fmt.Println(
		appConfig.KafkaHost,
		appConfig.KafkaPort,
		appConfig.KafkaTopic,
		appConfig.EventsToProcessChannelBufferSize,
		appConfig.ProcessedEventsChannelBufferSize,
	)

	mailHandler = consumer.New(*appConfig, logger.Sugar())

	mailHandler.Start(ctx)

}

func Stop() {
	mailHandler.Stop()
	logger.Sugar().Info("app has stopped")
}
