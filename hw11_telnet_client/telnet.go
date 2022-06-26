package main

import (
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
	return c.conn.Close()
}

func (c *simpleClient) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("sending error: %w", err)
	}
	return nil
}

func (c *simpleClient) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("receiving error: %w", err)
	}
	return nil
}
