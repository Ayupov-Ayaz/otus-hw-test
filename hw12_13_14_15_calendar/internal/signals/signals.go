package signals

import (
	"context"
	"os/signal"
	"syscall"
)

func NotifyCtx() (context.Context, func()) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}
