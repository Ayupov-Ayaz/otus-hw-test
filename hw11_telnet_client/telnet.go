package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := fmt.Fprintln(c.conn, scanner.Text())
		if err != nil {
			log.Printf("Failed to write to socket: %v", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read from STDIN: %v", err)
		return err
	}

	return nil
}

func (c *telnetClient) Receive() error {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Failed to read from socket: %v", err)
				return err
			}
			return nil
		}

		fmt.Fprint(c.out, line)
	}
}
