package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/spf13/viper"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		// printVersion()
		return
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err.Error() + `
		Default path to config is /etc/calendar/<config name>.yml. Check it's existence!
		If you want use own configuration use calendar --config=<path to config>`)
		return
	}

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	config := cfg.NewConfig()
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(
		config.HTTP.Host,
		config.HTTP.Port,
		logg,
		calendar,
	)

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
