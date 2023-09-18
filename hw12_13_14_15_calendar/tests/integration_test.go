//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/sender"
	rpcserver "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/server"
	internalhttp "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/migrations"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CalendarSuite struct {
	suite.Suite
	server    *internalhttp.Server
	rpcServer *rpcserver.GRPCServer
	calendar  *app.App
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	logg      *logger.Logger
	config    *cfg.Config
}

func (cs *CalendarSuite) SetupSuite() {
	fmt.Println("Setup suite")
	file, err := os.Open("/etc/calendar/config.yml")
	cs.Require().NoError(err)

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	cs.config = cfg.NewConfig()
	err = viper.Unmarshal(cs.config)
	cs.Require().NoError(err)

	cs.logg = logger.New(cs.config.Logger.Level, os.Stdout)

	storage := sqlstorage.New(cs.config.Storage.DSN)
	cs.calendar = app.New(cs.logg, storage)

	cs.server = internalhttp.NewServer(
		cs.config.HTTP.Host,
		cs.config.HTTP.Port,
		cs.config.HTTP.ReadHeaderTimeout,
		cs.logg,
		cs.calendar,
	)

	cs.rpcServer = rpcserver.NewGRPCServer(cs.logg, cs.calendar)

	cs.ctx, cs.cancel = context.WithCancel(context.Background())
	// // Start http server
	cs.wg.Add(2)
	go func() {
		defer cs.wg.Done()
		cs.server.Start()
	}()

	// // Start grpc server
	go func() {
		defer cs.wg.Done()
		cs.rpcServer.Start(net.JoinHostPort(cs.config.RPC.Host, cs.config.RPC.Port))
	}()
}

func (cs *CalendarSuite) SetupTest() {
	fmt.Println("Setup test")
	migrations.Up("files")
	migrations.Up("inserting")
}

// For this test in database should be events with exeed notification time
// to scheduler can accept them all together.
func (cs *CalendarSuite) TestScheduler() {
	file, err := os.Open("/etc/calendar/schedulercfg.yml")
	cs.Require().NoError(err)

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	cfg := scheduler.NewSchedulerConfig()
	err = viper.Unmarshal(cfg)
	cs.Require().NoError(err)

	sched := scheduler.NewScheduler(*cfg)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.Start(ctx)
	}()
	time.Sleep(5*time.Second + time.Duration(cfg.Timeout)*time.Second)
	cs.Require().NoError(err, "error in start scheduler occurred")
	events, err := cs.calendar.NotScheduledEvents()
	cs.Require().NoError(err, "error in selection unscheduled events occurred")
	cs.Require().Empty(events)
	events, err = cs.calendar.Events()
	cs.Require().NoError(err, "error in selection events occurred")
	for _, e := range events {
		cs.Require().True(e.Scheduled, "Event isn't scheduled")
	}
	cancel()
	sched.Stop()
	wg.Wait()
}

func (cs *CalendarSuite) TestSender() {
	file, err := os.Open("/etc/calendar/schedulercfg.yml")
	cs.Require().NoError(err)

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	cfg := scheduler.NewSchedulerConfig()
	err = viper.Unmarshal(cfg)
	cs.Require().NoError(err)

	sched := scheduler.NewScheduler(*cfg)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(2)
	go func() {
		defer wg.Done()
		sched.Start(ctx)
	}()

	file, err = os.Open("/etc/calendar/sendercfg.yml")
	cs.Require().NoError(err)

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	senderCfg := sender.NewSenderConfig()
	err = viper.Unmarshal(senderCfg)
	cs.Require().NoError(err)
	sender := sender.NewSender(*senderCfg)

	go func() {
		defer wg.Done()
		sender.Start(ctx)
	}()
	time.Sleep(5*time.Second + time.Duration(cfg.Timeout)*time.Second)
	cs.Require().NotEmpty(sender.Message, "Sender didn't get data from rabbit queue")
	cancel()
	sched.Stop()
	sender.Stop()
	wg.Wait()
}

func (cs *CalendarSuite) TearDownTest() {
	migrations.Down("inserting")
	migrations.Down("files")
	fmt.Println("Tear Down test")
}

func (cs *CalendarSuite) TearDownSuite() {
	if err := cs.server.Stop(cs.ctx); err != nil {
		fmt.Println(err)
	}
	cs.rpcServer.Stop()
	cs.cancel()
	cs.wg.Wait()
	fmt.Println("Tear Down suite")
}

func TestCalendarService(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
