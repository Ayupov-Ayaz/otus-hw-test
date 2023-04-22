package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/app"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/memory"
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

	config, err := NewConfig()
	if err != nil {
		return err
	}

	logg := logger.New(config.Logger.Level)
	defer logg.Sync()

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
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

	if err := server.Start(ctx); err != nil {
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
