package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &simpleClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type simpleClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *simpleClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.conn = conn

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)

	return nil
}

func (c *simpleClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *simpleClient) Send() error {
	writer := bufio.NewWriter(c.conn)
	if _, err := io.Copy(writer, c.in); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("sending error: %w", err)
	}
	return nil
}

func (c *simpleClient) Receive() error {
	reader := bufio.NewReader(c.conn)
	if _, err := io.Copy(c.out, reader); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("receiving error: %w", err)
	}
	return nil
}
