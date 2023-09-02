package logger

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	fileName := "log.txt"
	t.Run("info", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("info", logFile)
		time := time.Now().UTC().Format("01-12-2023 11:11")
		l.Info("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "[info] "+time+" Test text message\n", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
	t.Run("error", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("error", logFile)
		time := time.Now().UTC().Format("01-12-2023 11:11")
		l.Error("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "[error] "+time+" Test text message\n", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
	t.Run("warn", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("warn", logFile)
		time := time.Now().UTC().Format("01-12-2023 11:11")
		l.Warn("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "[warn] "+time+" Test text message\n", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
	t.Run("debug", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("debug", logFile)
		time := time.Now().UTC().Format("01-12-2023 11:11")
		l.Debug("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "[debug] "+time+" Test text message\n", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
}

func TestInvalid(t *testing.T) {
	fileName := "log.txt"
	t.Run("info", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("info", logFile)
		l.Debug("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
	t.Run("error", func(t *testing.T) {
		logFile, err := os.Create(fileName)
		if err != nil {
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		l := New("error", logFile)
		l.Info("Test text message")
		logFile.Seek(0, io.SeekStart)
		result, err := io.ReadAll(logFile)
		if err != nil {
			logFile.Close()
			os.Remove(fileName)
			t.Fatalf("%s\n", err.Error())
		}
		logFile.Close()
		require.Equal(t, "", string(result), "Not equal! Actual result is %q", result)
		os.Remove(fileName)
	})
}
