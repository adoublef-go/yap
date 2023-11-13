package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adoublef/yap/cmd/yap/server"
	"github.com/adoublef/yap/errgroup"
	"github.com/adoublef/yap/nats"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	// prepare server options
	fs := flag.NewFlagSet("", flag.ExitOnError)
	// TODO set usage
	opts, err := server.Parse(fs, os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	// run server
	if err := run(ctx, opts); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, opts *server.Options) (err error) {
	// prepare nats server
	ns, err := nats.NewServer(ctx)
	if err != nil {
		return fmt.Errorf("setup nats-server: %w", err)
	}
	ns.WaitForServer()
	defer ns.Close()
	// prepare nats connection
	nc, err := nats.Connect(ns)
	if err != nil {
		return fmt.Errorf("setup nats connection: %w", err)
	}
	defer nc.Close()
	// prepare http server
	s, err := server.NewServer(ctx, nc, opts)
	if err != nil {
		return fmt.Errorf("setup server: %w", err)
	}
	g := errgroup.New(ctx)
	g.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return s.Shutdown(context.Background())
	})
	return g.Wait()
	// errCh := make(chan error, 1)
	// go func() {
	// 	errCh <- s.ListenAndServe()
	// }()
	// select {
	// case err := <-errCh:
	// 	return err
	// case <-ctx.Done():
	// 	return s.Shutdown(context.Background())
	// }
}
