package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KafkaHost                        string `split_words:"true"`
	KafkaPort                        string `split_words:"true"`
	KafkaTopic                       string `split_words:"true"`
	WorkersCount                     int    `split_words:"true"`
	EventsToProcessChannelBufferSize int    `split_words:"true"`
	ProcessedEventsChannelBufferSize int    `split_words:"true"`
	MailRateLimitSec                 int    `split_words:"true"`
}

func NewConfig() (*Config, error) {
	var s Config

	err := envconfig.Process("", &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
