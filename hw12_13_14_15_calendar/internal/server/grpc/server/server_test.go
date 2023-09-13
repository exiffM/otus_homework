package rpcserver

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"hw12_13_14_15_calendar/internal/app"
	"hw12_13_14_15_calendar/internal/logger"
	rpcclient "hw12_13_14_15_calendar/internal/server/grpc/client"
	eventrpcapi "hw12_13_14_15_calendar/internal/server/grpc/pb"
	sqlstorage "hw12_13_14_15_calendar/internal/storage/sql"
	"hw12_13_14_15_calendar/migrations"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var dsn = "user=igor dbname=calendardb password=igor"

func TestIntegration(t *testing.T) {
	migrations.Up("files")
	gdsn := "localhost:5000"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := logger.New("info", os.Stdout)
	storage := sqlstorage.New(dsn)
	calendar := app.New(log, storage)
	gserver := NewGRPCServer(log, calendar)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info("GracefulStop")
		gserver.GracefulStop()
	}()

	go func() {
		defer wg.Done()
		log.Info("server listening at %v", "localhost:5000")
		gserver.Start(gdsn)
	}()
	time.Sleep(5 * time.Second)

	gclient := rpcclient.Client{}
	gclient.Connect(gdsn)
	log.Info("Client was created and connected")
	// Create
	gevent := eventrpcapi.Event{
		Tittle:      "Tittle of event 1",
		Start:       timestamppb.Now(),
		Duration:    1000000000,
		Description: "Description of event 1",
	}
	firstCreated, err := gclient.CreateEvent(ctx, &gevent)
	if err != nil {
		log.Info(err.Error())
		return
	}
	require.Equal(t, int32(1), firstCreated.Id,
		"Not equal! Actual Id of created event is: %v", firstCreated.Id)
	gevent = eventrpcapi.Event{
		Tittle:      "Tittle of event 2",
		Start:       timestamppb.Now(),
		Duration:    1000000000,
		Description: "Description of event 2",
	}
	createdGevent, err := gclient.CreateEvent(ctx, &gevent)
	if err != nil {
		log.Info(err.Error())
		return
	}
	require.Equal(t, int32(2), createdGevent.Id,
		"Not equal! Actual Id of created event is: %v", createdGevent.Id)
	// Update
	gevent = eventrpcapi.Event{
		Id:          1,
		Tittle:      "Update tittle if event 1",
		Start:       timestamppb.Now(),
		Duration:    1000000000,
		Description: "Update description of event 1",
	}
	updatedGevent, err := gclient.UpdateEvent(ctx, &gevent)
	if err != nil {
		log.Info(err.Error())
		return
	}
	require.Equal(t, gevent.Tittle, updatedGevent.Tittle,
		"Not equal! Actual titt;e is: %v", updatedGevent.Tittle)
	require.Equal(t, gevent.Description, updatedGevent.Description,
		"Not equal! Actual description is: %v", updatedGevent.Description)
	// Delete
	id := eventrpcapi.Id{Id: 1}
	id.Id = 1
	_, err = gclient.DeleteEvent(ctx, &id)
	if err != nil {
		log.Info(err.Error())
		return
	}
	deletedEvent, err := gclient.SelectEvent(ctx, &id)
	require.NotNil(t, err, "Error should be not nil!")
	require.Nil(t, deletedEvent, "Event is not nil!")
	// Select
	id.Id = 2
	selectedEvent, err := gclient.SelectEvent(ctx, &id)
	if err != nil {
		log.Info(err.Error())
		return
	}
	require.Equal(t, createdGevent.Tittle, selectedEvent.Tittle,
		"Not equal! Actual event's tittle is: %v", selectedEvent.Tittle)
	require.Equal(t, createdGevent.Duration, selectedEvent.Duration,
		"Not equal! Actual event's Duration is: %v", selectedEvent.Duration)
	// Show list
	events, err := gclient.Events(ctx)
	if err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("%+v", events)
	log.Info("All test passed successfully")
	gclient.Close()
	migrations.Down("files")
	cancel()
	wg.Wait()
}
