package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second, "timeout 2s")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Host or port arguments must be provided.")
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout duration] host port\n", os.Args[0])
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to establish connection : %s", err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close connection: %s", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT)
	defer stop()

	go receiver(cancel, client)
	go sender(cancel, client)

	<-ctx.Done()
}

func sender(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "sender error: %s", err)
		return
	}
}

func receiver(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Receive(); err != nil {
		fmt.Fprintf(os.Stderr, "receiver error: %s", err)
	}
}
