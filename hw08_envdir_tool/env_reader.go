package main

import (
	"bufio"
	"bytes"
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

func clearTerminalZeros(str string) string {
	return strings.Split(string(bytes.Replace([]byte(str), []byte{0}, []byte("\n"), -1)), "\n")[0]
}

func remove(str, remove string) string {
	return strings.Replace(str, remove, "", -1)
}

func clearString(str string) string {
	return clearTerminalZeros(remove(remove(remove(str, " "), "\t"), "\""))
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
