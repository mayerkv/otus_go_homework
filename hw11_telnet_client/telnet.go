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

var (
	ErrConnectionError = errors.New("connection error")
	ErrSendingError    = errors.New("sending error")
	ErrReceivingError  = errors.New("receiving error")
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
		return fmt.Errorf("%w: %s", ErrConnectionError, err.Error())
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
		return fmt.Errorf("%w: %s", ErrSendingError, err)
	}
	return nil
}

func (c *simpleClient) Receive() error {
	reader := bufio.NewReader(c.conn)
	if _, err := io.Copy(c.out, reader); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: %s", ErrReceivingError, err)
	}
	return nil
}
