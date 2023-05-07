package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/cmd/calendar/internal"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/app"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/logger"
	internalhttp "github.com/ayupov-ayaz/otus-wh-test/hw12/internal/server/http"
)

var configFile string

const defaultConfigFile = "/etc/calendar/config.toml"

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "Path to configuration file")
}

func run() error {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return nil
	}

	config, err := internal.NewConfig(configFile)
	if err != nil {
		return err
	}

	logg := logger.New(config.Logger.Level)
	defer func() {
		_ = logg.Sync()
	}()

	logg.Info("using config file", zap.String("path", configFile))
	logg.Info("using storage", zap.String("driver", config.Storage.Driver))

	storage, err := NewStorage(config.Storage)
	if err != nil {
		return err
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, release)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx, config.HTTP.Addr()); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
