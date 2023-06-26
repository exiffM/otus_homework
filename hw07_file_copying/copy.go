package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrEmptyInputFile        = errors.New("empty file")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrPermissionDenied      = errors.New("permission denied")
	ErrLimitIsInvalid        = errors.New("invalid limit value")
	ErrOffsetIsInvalid       = errors.New("invalid offset value")
	ErrEmptyOutputFilePath   = errors.New("empty output file path")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	switch {
	case err != nil || info.IsDir():
		return ErrUnsupportedFile
	case info.Size() == 0:
		return ErrEmptyInputFile
	case toPath == "":
		return ErrEmptyOutputFilePath
	case info.Mode() < 0o644:
		return ErrPermissionDenied
	case limit < 0:
		return ErrLimitIsInvalid
	case offset < 0:
		return ErrOffsetIsInvalid
	case offset > info.Size():
		return ErrOffsetExceedsFileSize
	default:
		if limit == 0 || limit > info.Size() {
			limit = info.Size()
		} else if offset > 0 && offset+limit > info.Size() {
			limit = info.Size() - offset
		}
		inputFile, err := os.Open(fromPath)
		defer inputFile.Close()
		if err != nil {
			return err
		}

		progress := pb.Simple.Start64(limit)
		progressReader := progress.NewProxyReader(inputFile)

		outputFile, err := os.Create(toPath)
		defer outputFile.Close()
		if err != nil {
			return err
		}

		inputFile.Seek(offset, io.SeekStart)

		n, err := io.CopyN(outputFile, progressReader, limit)
		_ = n
		return err
	}
}
