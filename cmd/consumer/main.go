package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/pflag"

	"github.com/alvarocabanas/cart/cmd/consumer/bootstrap"
	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

const serviceName = "cart-consumer"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	cfg := parseConfig()

	err := bootstrap.InitTraceExporter(cfg)
	if err != nil {
		log.Fatalf("error initializing trace exporter, %v", err)
	}

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
		fmt.Println(err)
		return err
	})

	g.Go(func() error {
		fmt.Printf("Launching Server")
		return metricServer.ListenAndServe()
	})

	err = g.Wait()
}

func parseConfig() bootstrap.Config {
	var cfg bootstrap.Config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	flagSet := pflag.NewFlagSet(serviceName, pflag.ContinueOnError)
	if err := gpflag.ParseTo(&cfg, flagSet, sflags.FlagDivider("."), sflags.FlagTag("mapstructure")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
	viper.SetConfigName(serviceName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s \n", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.ServiceName = serviceName
	return cfg
}
