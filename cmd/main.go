package main

import (
	"context"
	"mail/internal/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	go application.Start(ctx)
	<-ctx.Done()
	application.Stop()
}
