package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var stopFunc func()
var stopContext context.Context

func Init() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	stopFunc = stop
	stopContext = ctx

	return ctx
}

func Stop() {
	stopFunc()
}

func GetStopContext() context.Context {
	return stopContext
}
