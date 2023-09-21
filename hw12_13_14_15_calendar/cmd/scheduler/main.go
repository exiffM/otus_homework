package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/spf13/viper"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/schedulercfg.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err.Error() + `
		Default path to config is /etc/calendar/<config name>.yml. Check it's existence!
		If you want use own configuration use calendar --config=<path to config>`)
		return
	}

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	cfg := scheduler.NewSchedulerConfig()
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}

	wg := sync.WaitGroup{}
	onErrOnce := sync.Once{}
	wg.Add(2)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt,
		os.Kill)
	defer cancel()

	host := os.Getenv("CALENDAR_HOST")
	if host != "" {
		hp := strings.Split(cfg.Source.ConnectionString, ":")
		cfg.Source.ConnectionString = strings.Join([]string{host, hp[1]}, ":")
	}
	host = os.Getenv("RABBIT_HOST")
	if host != "" {
		cfg.Target.ConnectionString = host
	}
	scheduler := scheduler.NewScheduler(*cfg)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := scheduler.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := scheduler.Start(ctx); err != nil {
			fmt.Println(err)
			onErrOnce.Do(cancel)
		}
	}()
	wg.Wait()
}
