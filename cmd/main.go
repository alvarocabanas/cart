package main

import (
	"cart/cmd/bootstrap"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// In future iterations of the application all the config will come from environment variables
const serverPort = ":8888"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	cfg := bootstrap.Config{ServerPort: serverPort}

	server, err := bootstrap.InitializeServer(gctx, cfg)
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
