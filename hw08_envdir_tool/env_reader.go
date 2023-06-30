package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func replaceZeroTerminatedBytes(str string) string {
	byteArray := []byte(str)
	if bytes.Contains(byteArray, []byte("\x00")) {
		str = string(bytes.Join(bytes.Split(byteArray, []byte("\x00")), []byte("\n")))
	}
	str = strings.TrimRight(str, " ")

	return str
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(map[string]EnvValue)
	entryList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, path := range entryList {
		file, err := os.Open(dir + "/" + path.Name())
		if err != nil {
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if stat.Size() == 0 {
			env[path.Name()] = EnvValue{"", true}
		} else {
			fileScanner := bufio.NewScanner(file)
			fileScanner.Split(bufio.ScanLines)
			fileScanner.Scan()
			env[path.Name()] = EnvValue{replaceZeroTerminatedBytes(fileScanner.Text()), false}
		}

		file.Close()
	}
	return env, nil
}
