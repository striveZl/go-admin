package util

import (
	"context"
	"go-admin/pkg/logging"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, handler func(ctx context.Context) (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sc)

	cleanFn, err := handler(ctx)
	if err != nil {
		return err
	}

EXIT:
	for {
		select {
		case sig := <-sc:
			logging.Context(ctx).Info("Received signal", zap.String("signal", sig.String()))

			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				state = 0
				break EXIT
			case syscall.SIGHUP:
			default:
				break EXIT
			}
		case <-ctx.Done():
			logging.Context(ctx).Info("Context canceled, preparing to exit", zap.Error(ctx.Err()))
			state = 0
			break EXIT
		}
	}

	if cleanFn != nil {
		cleanFn()
	}
	logging.Context(ctx).Info("Server exit, bye...")
	time.Sleep(time.Millisecond * 100)
	os.Exit(state)
	return nil
}
