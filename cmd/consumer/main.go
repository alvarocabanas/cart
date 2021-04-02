package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alvarocabanas/cart/cmd/consumer/bootstrap"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	err := bootstrap.InitTraceExporter()
	if err != nil {
		panic(err)
	}

	cfg := bootstrap.Config{}
	consumer, err := bootstrap.InitializeConsumer(ctx, cfg)
	if err != nil {
		panic(err)
	}

	metricServer, err := bootstrap.InitMetricsServer(cfg)
	if err != nil {
		panic(err)
	}

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
			cancel()
		case <-ctx.Done():
			fmt.Printf("closing signal goroutine\n")
			metricServer.Shutdown(ctx)
			return ctx.Err()
		}

		return nil
	})

	g.Go(func() error {
		fmt.Printf("Launching Consumer")
		err := consumer.Start(gctx)
		fmt.Println("Consumer shutdown")
		return err
	})

	g.Go(func() error {
		fmt.Printf("Launching Server")
		return metricServer.ListenAndServe()
	})

	err = g.Wait()
}
