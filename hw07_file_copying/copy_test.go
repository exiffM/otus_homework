package main

import (
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("Unsupported file", func(t *testing.T) {
		err := Copy("testdata", "out.txt", 0, 0)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error is %v", err)
	})

	t.Run("Empty input file", func(t *testing.T) {
		err := Copy("/dev/random", "out.txt", 0, 0)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrEmptyInputFile), "actual error is %v", err)
	})

	t.Run("Empty output file path", func(t *testing.T) {
		err := Copy("testdata/input.txt", "", 0, 0)
		require.Truef(t, errors.Is(err, ErrEmptyOutputFilePath), "actual error is %v", err)
	})

	t.Run("No permission file", func(t *testing.T) {
		os.Chmod("testdata/forbiden.txt", 0o222)
		defer os.Chmod("testdata/forbiden.txt", 0o644)
		err := Copy("testdata/forbiden.txt", "out.txt", 0, 0)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrPermissionDenied), "actual error is %v", err)
	})

	t.Run("Invalid limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 0, -1)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrLimitIsInvalid), "actual error is %v", err)
	})

	t.Run("Invalid offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", -10, 1)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrOffsetIsInvalid), "actual error is %v", err)
	})

	t.Run("Offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 6618, 1)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error is %v", err)
	})

	t.Run("Sum of limit & offset greater than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 6000, 1000)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, nil), "actual error is %v", err)
	})

	t.Run("Successful work, full file", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, nil), "actual error is %v", err)
		f1, _ := os.Open("testdata/input.txt")
		defer f1.Close()
		f2, _ := os.Open("out.txt")
		defer f2.Close()
		defer os.Remove(f2.Name())

		h1 := sha256.New()
		h2 := sha256.New()
		if _, err := io.Copy(h1, f1); err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(h2, f2); err != nil {
			log.Fatal(err)
		}
		require.Equal(t, h1, h2, "Files' contents are different")
	})

	t.Run("Successful work, cropped file", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 6000, 1000)
		require.Truef(t, errors.Is(err, nil), "actual error is %v", err)
		f1, _ := os.Open("testdata/out_offset6000_limit1000.txt")
		defer f1.Close()
		f2, _ := os.Open("out.txt")
		defer f2.Close()
		defer os.Remove(f2.Name())

		h1 := sha256.New()
		h2 := sha256.New()
		if _, err := io.Copy(h1, f1); err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(h2, f2); err != nil {
			log.Fatal(err)
		}
		require.Equal(t, h1, h2, "Files' contents are different")
	})
}
