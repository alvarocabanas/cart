// +build wireinject

package bootstrap

import (
	"context"
	"net/http"

	"github.com/alvarocabanas/cart/internal/io/async"

	"github.com/alvarocabanas/cart/internal/messaging"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	cart "github.com/alvarocabanas/cart/internal"
	"github.com/alvarocabanas/cart/internal/metrics"
	"github.com/google/wire"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

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
	getMessageHandlerFunc,
	messaging.NewKafkaConsumer,
)

func InitializeConsumer(ctx context.Context, cfg Config) (*messaging.KafkaConsumer, error) {
	wire.Build(
		metrics.NewOpenCensusRecorder,
		wire.Bind(new(metrics.Recorder), new(metrics.OpenCensusRecorder)),
		messagingSet,
	)
	return &messaging.KafkaConsumer{}, nil
}

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
