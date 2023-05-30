package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
)

var (
	host    string
	port    string
	timeout = flag.Duration("timeout", 10*time.Second, "connection timeout")
)

func parseArgs() error {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		return fmt.Errorf("invalid number of arguments: %d", len(args))
	}

	host = args[0]
	port = args[1]

	return nil
}

func main() {
	if err := parseArgs(); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), *timeout)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	go readFromSocket(reader, NewStdoutPrinter())

	writeToSocket(ctx, conn)
}
