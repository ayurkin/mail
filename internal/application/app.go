package application

import (
	"context"
	consumer "mail/internal/adapters/kafka"
	"mail/internal/config"

	"go.uber.org/zap"
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

	mailHandler = consumer.New(appConfig, logger.Sugar())

	mailHandler.Start(ctx)
}

func Stop() {
	mailHandler.Stop()
	logger.Sugar().Info("app has stopped")
}
