package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-sigChan
		cancel()

		<-sigChan
		os.Exit(1)
	}()

	return ctx
}
