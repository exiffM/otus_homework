package sqlstorage

import (
	"context"
	"testing"
	"time"

	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

// var dsn = "user=igor password=igor host='127.0.0.1' database=calendardb"

var dsn = "user=igor dbname=calendardb password=igor"

func TestStorage(t *testing.T) {
	storage := New(dsn)
	ctx := context.TODO()
	loc, err := time.LoadLocation("Europe/Moscow")
	_ = loc // hotfix for tests (TODO: accept location on application level)
	_ = err

	t.Run("connect to database", func(t *testing.T) {
		err := storage.Connect(ctx)
		require.Nil(t, err, "Failed! Actual error is not nil!")
	})

	t.Run("create/select event", func(t *testing.T) {
		event := mdl.Event{
			ID:               1,
			Title:            "First event",
			Start:            time.Now().Add(time.Minute * 30).Round(time.Minute).In(loc),
			Duration:         int(time.Second * 10),
			Description:      "My first event in this calendar",
			NotificationTime: int(time.Minute * 15),
		}
		err := storage.CreateEvent(event)
		require.Nil(t, err, "Failed! Actual error is not nil!")
		nEvent, err := storage.SelectEvent(1)
		require.Nil(t, err, "Failed! Actual error is not nil!")
		require.Equal(t, event, nEvent, "Failed! Events aren't equal!")
	})

	t.Run("update event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               1,
			Title:            "Updated event",
			Start:            time.Now().Add(time.Minute * 30).Round(time.Minute).In(loc),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		err := storage.UpdateEvent(event)
		require.Nil(t, err, "Failed! Error is not nil!")

		nEvent, _ := storage.SelectEvent(1)
		require.Equal(t, event.Title, nEvent.Title,
			"Failed! Tittle of event hasn't been updated")
		require.Equal(t, event.Description, nEvent.Description,
			"Failed! Description of event hasn't been updated")
	})

	t.Run("delete event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               2, // has no effects
			Title:            "New event",
			Start:            time.Now().Add(time.Minute * 30).Round(time.Minute).In(loc),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		err := storage.CreateEvent(event)
		require.Nil(t, err, "Failed! Error is not nil")

		err = storage.DeleteEvent(1)
		require.Nil(t, err, "Failed! Error is not nil")
		ev, err := storage.SelectEvent(1)
		require.NotNil(t, err, "Failed! Error hasn't been ocure")
		// require.ErrorIs(t, err, errExist, "Failed! Actual error is %q", err)
		require.Equal(t, mdl.Event{}, ev, "Failed! Event isn't empty")
	})

	t.Run("list events", func(t *testing.T) {
		events, err := storage.Events()
		require.Nil(t, err, "Failed! Error is not nil")
		expectedEvent, err := storage.SelectEvent(2)
		require.Nil(t, err, "Failed! Error is not nil")
		require.Equal(t, []mdl.Event{expectedEvent},
			events, "Failed! Event slices are not equal!")
	})
}
