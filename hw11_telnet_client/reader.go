package main

import (
	"fmt"
	"io"
)

const delim = '\n'

var errEndOfInput = fmt.Errorf("end of input")

type StringReader interface {
	ReadString(delim byte) (string, error)
}

type Printer interface {
	Print(str string)
	Err(err error)
}

type StdoutPrinter struct{}

func NewStdoutPrinter() *StdoutPrinter {
	return &StdoutPrinter{}
}

func (p *StdoutPrinter) Print(str string) {
	fmt.Print(str)
}

func (p *StdoutPrinter) Err(err error) {
	fmt.Printf("Error: %v\n", err)
}

func read(reader StringReader) (line string, err error) {
	line, err = reader.ReadString(delim)
	if err != nil && err == io.EOF {
		err = errEndOfInput
	}

	return line, err
}

func readFromSocket(reader StringReader, printer Printer) {
	var (
		line string
		err  error
	)

	for {
		line, err = read(reader)
		if err != nil {
			if err == errEndOfInput {
				return
			}

			printer.Err(fmt.Errorf("Failed to read from socket: %v\n", err))
			continue
		}

		printer.Print(line)
	}
}
