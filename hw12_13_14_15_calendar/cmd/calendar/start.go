package main

import (
	"context"
	"fmt"
	app "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/server/grpc/api"
	internalhttp "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/server/http"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
)

type Shutdown func(ctx context.Context) error

func shutdown(ctx context.Context, logg *zap.Logger, shutdownHTTP, ShutdownGRPC Shutdown) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := shutdownHTTP(ctx); err != nil {
			logg.Error("failed to stop http server", zap.Error(err))
		}
		wg.Done()
	}()

	go func() {
		if err := ShutdownGRPC(ctx); err != nil {
			logg.Error("failed to stop grpc server", zap.Error(err))
		}
		wg.Done()
	}()

	wg.Wait()
}

func startHTTP(ctx context.Context, calendar *app.EventUseCase, logg *zap.Logger, addr string) (Shutdown, error) {
	server := internalhttp.NewServer(logg, release)
	server.Register(internalhttp.NewEventHandlers(logg, calendar))

	go func() {
		if err := server.Start(ctx, addr); err != nil {
			logg.Error("failed to start http server: ", zap.Error(err))
			os.Exit(1)
		}
	}()

	return server.Stop, nil
}

func startGRPC(calendar *app.EventUseCase, logg *zap.Logger, addr string) (Shutdown, error) {
	server := grpc.NewServer(grpc.StreamInterceptor(api.LoggerInterceptor(logg)))
	handler := api.NewEventHandler(calendar, logg)
	api.RegisterEventServiceServer(server, handler)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			logg.Error("failed to start grpc server: ", zap.Error(err))
			os.Exit(1)
		}
	}()

	stop := func(_ context.Context) error {
		server.GracefulStop()
		return nil
	}

	return stop, nil
}
