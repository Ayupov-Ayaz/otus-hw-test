package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	address string
	timeout = flag.Duration("timeout", 10*time.Second, "connection timeout")
)

var (
	errInvalidNumberOfArguments = errors.New("invalid number of arguments")
	errHostIsEmpty              = errors.New("host is empty")
	errPortIsEmpty              = errors.New("port is empty")
)

func parseArgs(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("%d: %w", len(args), errInvalidNumberOfArguments)
	}

	host := args[0]
	port := args[1]

	if host == "" {
		return errHostIsEmpty
	}

	if port == "" {
		return errPortIsEmpty
	}

	address = net.JoinHostPort(host, port)

	return nil
}

func main() {
	flag.Parse()
	if err := parseArgs(flag.Args()); err != nil {
		log.Fatal(fmt.Errorf("parse flags failed: %w", err))
	}

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	go func() {
		err := client.Receive()
		if err != nil {
			log.Printf("error receiving data: %v", err)
		}
	}()

	err = client.Send()
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}

	// Ожидание сигнала SIGINT для завершения программы
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	log.Println("Program terminated")
}
