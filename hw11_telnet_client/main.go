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
	errInvalidNumberOfArguments = errors.New("invalid number of arguments")
	errHostIsEmpty              = errors.New("host is empty")
	errPortIsEmpty              = errors.New("port is empty")
)

func parseArgs(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("%d: %w", len(args), errInvalidNumberOfArguments)
	}

	host := args[0]
	port := args[1]

	if host == "" {
		return "", errHostIsEmpty
	}

	if port == "" {
		return "", errPortIsEmpty
	}

	address := net.JoinHostPort(host, port)

	return address, nil
}

func getTimeout() time.Duration {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	return *timeout
}

func run(address string, timeout time.Duration) error {
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		return err
	}

	defer func() {
		log.Println(client.Close())
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			log.Println(fmt.Errorf("failed to receive: %w", err))
		}
	}()

	if err := client.Send(); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// Ожидание сигнала SIGINT для завершения программы
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	log.Println("Program terminated")

	return nil
}

func main() {
	timeout := getTimeout()
	address, err := parseArgs(flag.Args())
	if err != nil {
		log.Fatal(fmt.Errorf("parse flags failed: %w", err))
	}

	if err := run(address, timeout); err != nil {
		log.Fatal(err)
	}
}
