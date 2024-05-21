package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/app"
	config2 "github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/storage/database"
	memorystorage "github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/pkg/shortcuts"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./build/local/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx := context.Background()

	config, err := config2.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(config.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	var eventStorage storage.IStorage
	if config.Storage.InDatabase() {
		db := database.New(&config.Database)
		err := db.Connect(ctx)
		shortcuts.FatalIfErr(err)

		defer func() {
			err := db.Connect(ctx)
			if err != nil {
				logg.Error(err)
			}
		}()

		eventStorage = sqlstorage.NewEventStorage(db, logg)
	} else {
		eventStorage = memorystorage.New()
	}

	calendar := app.New(logg, eventStorage)

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
}
