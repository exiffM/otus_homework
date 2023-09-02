package memorystorage

import (
	"testing"
	"time"

	mdl "hw12_13_14_15_calendar/internal/storage"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()
	t.Run("connect to database", func(t *testing.T) {
		err := storage.Connect()
		require.Nil(t, err, "Failed! Actual error is not nil!")
	})

	t.Run("create event", func(t *testing.T) {
		event := mdl.Event{
			Title:            "First event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "My first event in this calendar",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.CreateEvent(event)
		require.Nil(t, err, "Failed! Actual error is not nil!")
		require.Equal(t, event, storage.data[0], "Failed! Events aren't equal!")
	})

	t.Run("select event", func(t *testing.T) {
		event, err := storage.SelectEvent(0)
		require.Nil(t, err, "Failed! Actual error is not nil!")
		require.Equal(t, storage.data[0], event, "Failed! Events aren't equal!")

		event, err = storage.SelectEvent(1)
		require.ErrorIs(t, err, errExist, "Failed! Actual error is %q", err)
		require.Empty(t, event, "Failed! Event isn't empty")
	})

	t.Run("update event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               0,
			Title:            "Updated event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.UpdateEvent(event)
		require.Nil(t, err, "Failed! Error is not nil!")

		event.ID = 1
		_, err = storage.UpdateEvent(event)
		require.ErrorIs(t, err, errExist, "Failed! Actual error is %q", err)
	})

	t.Run("delete event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               1,
			Title:            "New event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.CreateEvent(event)
		require.Nil(t, err, "Failed! Error is not nil")

		err = storage.DeleteEvent(0)
		require.Nil(t, err, "Failed! Error is not nil")
		ev, err := storage.SelectEvent(0)
		require.ErrorIs(t, err, errExist, "Failed! Actual error is %q", err)
		require.Empty(t, ev, "Failed! Event isn't empty")
	})

	t.Run("list events", func(t *testing.T) {
		events, err := storage.Events()

		require.Nil(t, err, "Failed! Error is not nil")
		require.Equal(t, []mdl.Event{storage.data[1]},
			events, "Failed! Event slices are not equal!")
	})
}

func TestStorageInvalid(t *testing.T) {
	storage := New()

	t.Run("create event", func(t *testing.T) {
		event := mdl.Event{
			Title:            "First event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "My first event in this calendar",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.CreateEvent(event)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
	})

	t.Run("select event", func(t *testing.T) {
		event, err := storage.SelectEvent(0)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
		require.Empty(t, event, "Failed! Event isn't empty")

		event, err = storage.SelectEvent(1)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
		require.Empty(t, event, "Failed! Event isn't empty")
	})

	t.Run("update event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               0,
			Title:            "Updated event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.UpdateEvent(event)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)

		event.ID = 1
		_, err = storage.UpdateEvent(event)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
	})

	t.Run("delete event by id", func(t *testing.T) {
		event := mdl.Event{
			ID:               1,
			Title:            "New event",
			Start:            time.Now().Add(time.Minute * 30),
			Duration:         int(time.Second * 10),
			Description:      "Other description of event",
			NotificationTime: int(time.Minute * 15),
		}
		_, err := storage.CreateEvent(event)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)

		err = storage.DeleteEvent(0)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)

		ev, err := storage.SelectEvent(0)
		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
		require.Empty(t, ev, "Failed! Event isn't empty")
	})

	t.Run("list events", func(t *testing.T) {
		events, err := storage.Events()

		require.ErrorIs(t, err, errCreate, "Failed! Actual error is %q", err)
		require.Empty(t, events, "Failed! Event isn't empty")
	})
}
