package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/alvarocabanas/cart/cmd/consumer/bootstrap"
	"golang.org/x/sync/errgroup"
)

const serviceName = "cart-consumer"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	var cfg bootstrap.Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	err = bootstrap.InitTraceExporter(cfg)
	if err != nil {
		log.Fatalf("error initializing trace exporter, %v", err)
	}
	cfg.ServiceName = serviceName

	consumer, err := bootstrap.InitializeConsumer(ctx, cfg)
	if err != nil {
		log.Fatalf("error initializing consumer, %v", err)
	}

	metricServer, err := bootstrap.InitMetricsServer(cfg)
	if err != nil {
		log.Fatalf("error initializing metrics server, %v", err)
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
