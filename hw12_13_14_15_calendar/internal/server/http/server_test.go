package internalhttp

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/logger"
	api "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http/api/default"
	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/require"
)

// func init() {
// 	ctx = context.Background()
// 	log = logger.New("info", os.Stdin)
// 	source = sqlstorage.New(dsn)
// 	application = app.New(log, source)
// 	httpServer = NewServer(host, port, log, application)

//		wg = &sync.WaitGroup{}
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			httpServer.Start(ctx)
//		}()
//		time.Sleep(5 * time.Second)
//		client = &http.Client{}
//	}

func TestComplex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dsn := "user=igor dbname=calendardb password=igor"
	host := "localhost"
	port := "1235"
	log := logger.New("info", os.Stdin)
	source := sqlstorage.New(dsn)
	application := app.New(log, source)
	httpServer := NewServer(host, port, log, application)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(ctx)
	}()
	defer func() {
		wg.Done()
		httpServer.Stop(ctx)
	}()
	time.Sleep(5 * time.Second)
	client := &http.Client{}
	t.Run("Create", func(t *testing.T) {
		sb := strings.Builder{}
		sb.WriteString("http://")
		sb.WriteString(net.JoinHostPort(host, port))
		sb.WriteString("/api/calendar/event")
		requestPath := sb.String()
		requestData := `{"Title": "Tittle of event 1",
			"Start": "2023-08-21T19:08:42+03:00",
			"Description": "Description of event 1"}`
		body := strings.NewReader(requestData)
		request, err := http.NewRequestWithContext(ctx, "POST", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		defResponse := api.NewEventResponse("", "", mdl.Event{})
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, 1, defResponse.Data.ID,
			"Invalid ID! Actual is: %q", defResponse.Data.ID)
		require.Equal(t, "Tittle of event 1", defResponse.Data.Title,
			"Invalid tittle! Actual is: %q", defResponse.Data.Title)
		require.Equal(t, "Description of event 1", defResponse.Data.Description,
			"Invalid description! Actual is: %q", defResponse.Data.Description)
		response.Body.Close()

		requestData = `{"Title": "Tittle of event 2",
			"Start": "2023-08-21T19:30:00+03:00",
			"Description": "Description of event 2"}`
		body = strings.NewReader(requestData)
		request, err = http.NewRequestWithContext(ctx, "POST", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err = client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, 2, defResponse.Data.ID,
			"Invalid ID! Actual is: %q", defResponse.Data.ID)
		require.Equal(t, "Tittle of event 2", defResponse.Data.Title,
			"Invalid tittle! Actual is: %q", defResponse.Data.Title)
		require.Equal(t, "Description of event 2", defResponse.Data.Description,
			"Invalid description! Actual is: %q", defResponse.Data.Description)
		response.Body.Close()

		requestData = `{"Title": "Tittle of event 3",
			"Start": "2023-08-26T14:35:00+03:00",
			"Description": "Description of event 3"}`
		body = strings.NewReader(requestData)
		request, err = http.NewRequestWithContext(ctx, "POST", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err = client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, 3, defResponse.Data.ID,
			"Invalid ID! Actual is: %q", defResponse.Data.ID)
		require.Equal(t, "Tittle of event 3", defResponse.Data.Title,
			"Invalid tittle! Actual is: %q", defResponse.Data.Title)
		require.Equal(t, "Description of event 3", defResponse.Data.Description,
			"Invalid description! Actual is: %q", defResponse.Data.Description)
		response.Body.Close()
	})

	t.Run("Select", func(t *testing.T) {
		sb := strings.Builder{}
		sb.WriteString("http://")
		sb.WriteString(net.JoinHostPort(host, port))
		sb.WriteString("/api/calendar/events/select/1")
		requestPath := sb.String()
		body := strings.NewReader("")
		request, err := http.NewRequestWithContext(ctx, "GET", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		defResponse := api.NewEventResponse("", "", mdl.Event{})
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, 1, defResponse.Data.ID,
			"Invalid ID! Actual is: %q", defResponse.Data.ID)
		require.Equal(t, "Tittle of event 1", defResponse.Data.Title,
			"Invalid tittle! Actual is: %q", defResponse.Data.Title)
		require.Equal(t, "Description of event 1", defResponse.Data.Description,
			"Invalid description! Actual is: %q", defResponse.Data.Description)
		response.Body.Close()
	})
	t.Run("Update", func(t *testing.T) {
		sb := strings.Builder{}
		sb.WriteString("http://")
		sb.WriteString(net.JoinHostPort(host, port))
		sb.WriteString("/api/calendar/events/update/1")
		requestPath := sb.String()
		requestData := `{"Title": "Tittle of changed event 1",
			"Start": "2023-08-21T19:08:42+03:00",
			"Description": "Changed description"}`
		body := strings.NewReader(requestData)
		request, err := http.NewRequestWithContext(ctx, "PUT", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		defResponse := api.NewEventResponse("", "", mdl.Event{})
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, 1, defResponse.Data.ID,
			"Invalid ID! Actual is: %q", defResponse.Data.ID)
		require.Equal(t, "Tittle of changed event 1", defResponse.Data.Title,
			"Invalid tittle! Actual is: %q", defResponse.Data.Title)
		require.Equal(t, "Changed description", defResponse.Data.Description,
			"Invalid description! Actual is: %q", defResponse.Data.Description)
		response.Body.Close()
	})

	t.Run("Delete", func(t *testing.T) {
		sb := strings.Builder{}
		sb.WriteString("http://")
		sb.WriteString(net.JoinHostPort(host, port))
		sb.WriteString("/api/calendar/events/delete/3")
		requestPath := sb.String()
		body := strings.NewReader("")
		request, err := http.NewRequestWithContext(ctx, "DELETE", requestPath, body)
		require.Nil(t, err, "Error in create request with context")
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		require.Nil(t, err, "Error in client.Do request")
		defResponse := api.NewEventResponse("", "", mdl.Event{})
		err = easyjson.UnmarshalFromReader(response.Body, defResponse)
		require.Nil(t, err, "Error in create while unmarshaling response to default")
		require.Equal(t, "You've successfully deleted event with id: 3", defResponse.Error,
			"Error in response! Actual response is: %q", defResponse.Error)
		response.Body.Close()
	})
}
