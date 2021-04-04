package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/alvarocabanas/cart/cmd/server/bootstrap"
	"golang.org/x/sync/errgroup"
)

const serviceName = "cart-server"

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
	cfg.ServiceName = serviceName

	err = bootstrap.InitTraceExporter(cfg)
	if err != nil {
		log.Fatalf("error initializing trace exporter, %v", err)
	}

	server, err := bootstrap.InitializeServer(gctx, cfg)
	if err != nil {
		log.Fatalf("error initializing server, %v", err)
	}

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
			cancel()
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return ctx.Err()
		}

		return nil
	})

	g.Go(func() error {
		fmt.Printf("Launching Server")
		err := server.ListenAndServe()
		fmt.Println("Server shutdown")
		return err
	})

	g.Go(func() error {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		fmt.Println("Shutting down Gracefully the server")
		return err
	})

	err = g.Wait()
}
