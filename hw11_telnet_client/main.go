package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	host := flag.String("host", "", "hostname or IP address")
	port := flag.Int("port", 0, "port number")
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if *host == "" || *port == 0 {
		log.Fatal("Please provide both host and port")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(*host, fmt.Sprintf("%d", *port)), *timeout)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	go readFromSocket(conn)
	writeToSocket(ctx, conn)
}

func readFromSocket(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from socket: %v", err)
			}
			return
		}
		fmt.Print(line)
	}
}

func writeToSocket(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			log.Printf("Failed to write to socket: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read from STDIN: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Println("Program terminated")
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for program termination")
	}
}
