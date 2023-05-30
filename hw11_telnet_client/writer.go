package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func writeToSocket(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := fmt.Fprintln(conn, scanner.Text())
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
