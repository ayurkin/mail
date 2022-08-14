package application

import (
	"context"
	consumer "mail/internal/adapters/kafka"
	"mail/internal/config"

	"go.uber.org/zap"
)

type App struct {
	logger      *zap.Logger
	mailHandler *consumer.Client
}

func Start(ctx context.Context, app *App) {
	logger, _ := zap.NewProduction()

	app.logger = logger
	appConfig, err := config.NewConfig()
	if err != nil {
		logger.Sugar().Fatalf("create config failed: %v", err)
	}

	app.mailHandler = consumer.New(appConfig, logger.Sugar())

	app.mailHandler.Start(ctx)
}

func Stop(app *App) {
	app.mailHandler.Stop()
	app.logger.Sugar().Info("app has stopped")
}
