package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
		fmt.Fprintf(os.Stderr, "Usage: %s [options] host port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to establish connection : %s", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
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
		log.Println(err)
		return
	}
}

func receiver(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Receive(); err != nil {
		log.Println(err)
	}
}
