// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package bootstrap

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/alvarocabanas/cart/internal"
	"github.com/alvarocabanas/cart/internal/io/async"
	"github.com/alvarocabanas/cart/internal/messaging"
	"github.com/alvarocabanas/cart/internal/metrics"
	"github.com/google/wire"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"net/http"
)

// Injectors from wire.go:

func InitializeConsumer(ctx context.Context, cfg Config) (*messaging.KafkaConsumer, error) {
	kafkaBrokers := getKafkaBrokers(cfg)
	kafkaConsumerGroupID := getKafkaConsumerGroupID(cfg)
	kafkaTopics := getKafkaTopics(cfg)
	openCensusRecorder, err := metrics.NewOpenCensusRecorder()
	if err != nil {
		return nil, err
	}
	handleFunc := getMessageHandlerFunc(cfg, openCensusRecorder)
	kafkaConsumer, err := messaging.NewKafkaConsumer(kafkaBrokers, kafkaConsumerGroupID, kafkaTopics, handleFunc)
	if err != nil {
		return nil, err
	}
	return kafkaConsumer, nil
}

// wire.go:

type Config struct {
	ServiceName string `mapstructure:"service_name"`
	Kafka       struct {
		Brokers []string `mapstructure:"brokers"`
		GroupID string   `mapstructure:"group_id"`
	} `mapstructure:"kafka"`
	MetricsServer struct {
		BindAddr string `mapstructure:"bind_addr"`
	} `mapstructure:"metrics_server"`
	JaegerTracingURL string `mapstructure:"jaeger_tracing_url"`
}

var messagingSet = wire.NewSet(
	getKafkaBrokers,
	getKafkaConsumerGroupID,
	getKafkaTopics,
	getMessageHandlerFunc, messaging.NewKafkaConsumer,
)

func InitMetricsServer(cfg Config) (*http.Server, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)

	mux := http.NewServeMux()
	metricsServer := &http.Server{Addr: cfg.MetricsServer.BindAddr, Handler: mux}
	mux.Handle("/metrics", exporter)

	return metricsServer, nil
}

func InitTraceExporter(cfg Config) error {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: cfg.JaegerTracingURL,
		Process: jaeger.Process{
			ServiceName: cfg.ServiceName,
		},
	})
	if err != nil {
		return err
	}
	defer exporter.Flush()
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}

func getKafkaBrokers(cfg Config) messaging.KafkaBrokers {
	return cfg.Kafka.Brokers
}

func getKafkaConsumerGroupID(cfg Config) messaging.KafkaConsumerGroupID {
	return messaging.KafkaConsumerGroupID(cfg.Kafka.GroupID)
}

func getKafkaTopics(cfg Config) messaging.KafkaTopics {
	return []string{cart.EventsTopic}
}

func getMessageHandlerFunc(cfg Config, metricsRecorder metrics.Recorder) messaging.HandleFunc {
	messageHandler := async.NewMessageHandler(metricsRecorder)
	return messageHandler.Handle
}
