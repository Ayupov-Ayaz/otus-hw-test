package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrValueInvalid = errors.New("env value invalid")

type Filter func(fileName string) bool

func equalFilter(exp string) Filter {
	return func(fileName string) bool {
		return fileName == exp
	}
}

func fileFormatFilter(exp string) Filter {
	return func(fileName string) bool {
		return strings.HasSuffix(fileName, exp)
	}
}

var fileFilters = []Filter{
	equalFilter(".DS_Store"),
	fileFormatFilter(".go"),
	fileFormatFilter(".sh"),
}

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
	return strings.Split(string(bytes.ReplaceAll([]byte(str), []byte{0}, []byte("\n"))), "\n")[0]
}

func remove(str, remove string) string {
	return strings.ReplaceAll(str, remove, "")
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

func skip(fileName string) bool {
	for _, filter := range fileFilters {
		if filter(fileName) {
			return true
		}
	}

	return false
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
			if skip(fileName) {
				continue
			}

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
