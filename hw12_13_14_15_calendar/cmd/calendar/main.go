package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"hw12_13_14_15_calendar/internal/app"
	cfg "hw12_13_14_15_calendar/internal/config"
	"hw12_13_14_15_calendar/internal/logger"
	rpcserver "hw12_13_14_15_calendar/internal/server/grpc/server"
	internalhttp "hw12_13_14_15_calendar/internal/server/http"
	memorystorage "hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// if flag.Arg(0) == "version" {
	// 	printVersion()
	// 	return
	// }

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
		config.HTTP.ReadHeaderTimeout,
		logg,
		calendar,
	)

	RPCServer := rpcserver.NewGRPCServer(logg, calendar)

	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	// 	<-c
	// 	cancel()
	// }()
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server.Start()
	})
	g.Go(func() error {
		return RPCServer.Start(net.JoinHostPort(config.RPC.Host, config.RPC.Port))
	})
	g.Go(func() error {
		<-gCtx.Done()
		RPCServer.GracefulStop()
		return server.Stop(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}

	stop()

	// ctx, cancel := signal.NotifyContext(context.Background(),
	// 	syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	// defer cancel()

	// wg := sync.WaitGroup{}
	// onErrorStopOnce := sync.Once{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	<-ctx.Done()

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()

	// 	if err := server.Stop(ctx); err != nil {
	// 		logg.Error("failed to stop http server: " + err.Error())
	// 	}

	// 	RPCServer.GracefulStop()
	// }()
	// // // Stop http server
	// // go func() {
	// // 	defer wg.Done()
	// // 	<-ctx.Done()
	// // 	fmt.Println("HTTP stopped")
	// // 	if err := server.Stop(ctx); err != nil {
	// // 		fmt.Println(err)
	// // 	}
	// // }()
	// // // Stop grpc server
	// // go func() {
	// // 	defer wg.Done()
	// // 	<-ctx.Done()
	// // 	fmt.Println("RPC stopped")
	// // 	RPCServer.Stop()
	// // }()

	// // Start http server
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := server.Start(); err != nil {
	// 		logg.Error("failed to start HTTP server: " + err.Error())
	// 		onErrorStopOnce.Do(cancel)
	// 	}
	// }()

	// // Start grpc server
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := RPCServer.Start(net.JoinHostPort(config.RPC.Host, config.RPC.Port)); err != nil {
	// 		logg.Error("failed to start RPC server: " + err.Error())
	// 		onErrorStopOnce.Do(cancel)
	// 	}
	// }()
	// logg.Error("Before wait")
	// wg.Wait()
	// logg.Error("After wait")
}
