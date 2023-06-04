package main

import (
	"flag"
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/signals"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/connect"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/validator"
	"log"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	app "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

const (
	defaultConfigFile = "/etc/calendar/config.toml"
	shutdownTimeout   = time.Second * 3
)

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

	storage, err := connect.NewStorage(config.Storage)
	if err != nil {
		return err
	}

	calendar := app.New(validator.New(), storage, logg)
	notifyCtx, cancel := signals.NotifyCtx()
	defer cancel()

	stopHTTP, err := startHTTP(notifyCtx, calendar, logg, config.HTTP.Addr())
	if err != nil {
		cancel()
		return fmt.Errorf("failed to start http server: %w", err)
	}

	stopGRPC, err := startGRPC(calendar, logg, config.GRPC.Addr())
	if err != nil {
		cancel()
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	logg.Info("calendar is running...")

	shutdown(notifyCtx, logg, stopHTTP, stopGRPC)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
