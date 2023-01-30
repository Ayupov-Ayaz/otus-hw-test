package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrValueInvalid = errors.New("env value invalid")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func NewEnv(value string) EnvValue {
	var needRemove bool
	if value == "" {
		needRemove = true
	}

	return EnvValue{Value: value, NeedRemove: needRemove}
}

func clearString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\\0", "\n", -1) // todo: ??????

	return str
}

func ReadFirstLineInFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	scanner := bufio.NewScanner(f)
	scanner.Scan() // читаем только первую строку

	resp := clearString(scanner.Text())

	return resp, scanner.Err()
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(name string) (Environment, error) {
	files, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}

	resp := make(Environment, len(files))

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			value, err := ReadFirstLineInFile(name + "/" + fileName)
			if err != nil {
				return nil, err
			}

			if strings.Contains(value, "=") {
				return nil, ErrValueInvalid
			}

			resp[fileName] = NewEnv(value)
		}

	}

	return resp, nil
}
