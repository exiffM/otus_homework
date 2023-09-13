package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	rpcserver "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/server"
	internalhttp "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/migrations"
	"github.com/spf13/viper"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
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

	storage := sqlstorage.New(config.Storage.DSN)
	calendar := app.New(logg, storage)

	if err := migrations.Up("files"); err != nil {
		log.Println("Unable to up migration in \"files\"")
	}
	if err := migrations.Up("inserting"); err != nil {
		log.Println("Unable to up migration in \"inserting\"")
	}
	// terr := migrations.Up("files")
	// _ = terr
	// terr = migrations.Up("preparedb")

	server := internalhttp.NewServer(
		config.HTTP.Host,
		config.HTTP.Port,
		config.HTTP.ReadHeaderTimeout,
		logg,
		calendar,
	)

	RPCServer := rpcserver.NewGRPCServer(logg, calendar)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	onErrorStopOnce := sync.Once{}
	// // Start http server
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil {
			logg.Error("failed to start HTTP server: " + err.Error())
			onErrorStopOnce.Do(cancel)
		}
	}()

	// // Start grpc server
	go func() {
		defer wg.Done()
		if err := RPCServer.Start(net.JoinHostPort(config.RPC.Host, config.RPC.Port)); err != nil {
			logg.Error("failed to start RPC server: " + err.Error())
			onErrorStopOnce.Do(cancel)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	if err := server.Stop(ctx); err != nil {
		fmt.Println(err)
	}
	RPCServer.Stop()
	onErrorStopOnce.Do(cancel)
	wg.Wait()
	// migrations.Down("preparedb")
	migrations.Down("inserting")
	migrations.Down("files")
}
