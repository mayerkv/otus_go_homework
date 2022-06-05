package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var ErrInvalidSymbol = errors.New("env name has invalid symbol")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(map[string]EnvValue)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		env, err := getEnvValue(dir, fileInfo)
		if err != nil {
			return nil, err
		}

		environment[fileInfo.Name()] = env
	}

	return environment, nil
}

func getEnvValue(dir string, info fs.FileInfo) (EnvValue, error) {
	if info.Size() == 0 {
		return EnvValue{
			Value:      "",
			NeedRemove: true,
		}, nil
	}

	if strings.Contains(info.Name(), "=") {
		return EnvValue{}, ErrInvalidSymbol
	}

	file, err := os.Open(path.Join(dir, info.Name()))
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return EnvValue{}, err
	}

	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))
	line = bytes.TrimRight(line, " ")
	line = bytes.TrimRight(line, "	")

	return EnvValue{
		Value:      string(line),
		NeedRemove: false,
	}, nil
}
