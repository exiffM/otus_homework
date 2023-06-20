package main

import (
	"errors"
	"io"
	"os"
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

// limit = 0 - whole file, limit > file size is ok
// offset <= fileSize
func Copy(fromPath, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	// check if file path is empty
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
		inputFile, err := os.Open(fromPath)
		defer inputFile.Close()
		if err != nil {
			return err
		}

		outputFile, err := os.Create(toPath)
		defer outputFile.Close()
		if err != nil {
			return err
		}

		inputFile.Seek(offset, io.SeekStart)

		if limit == 0 || limit > info.Size() {
			_, err := io.Copy(outputFile, inputFile)
			return err
		} else {
			_, err := io.CopyN(outputFile, inputFile, limit)
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}
