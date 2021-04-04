// +build wireinject

package bootstrap

import (
	"context"
	"net/http"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	cart "github.com/alvarocabanas/cart/internal"
	"github.com/alvarocabanas/cart/internal/creator"
	"github.com/alvarocabanas/cart/internal/getter"
	"github.com/alvarocabanas/cart/internal/io/rest"
	"github.com/alvarocabanas/cart/internal/messaging"
	"github.com/alvarocabanas/cart/internal/metrics"
	"github.com/alvarocabanas/cart/internal/storage"
	"github.com/google/wire"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type Config struct {
	ServiceName string `mapstructure:"service_name"`
	ServerPort  string `mapstructure:"server_port"`
	Kafka       struct {
		Brokers []string `mapstructure:"brokers"`
	} `mapstructure:"kafka"`
	JaegerTracingURL string `mapstructure:"jaeger_tracing_url"`
}

var appSet = wire.NewSet(
	getter.NewCartGetter,
	creator.NewCartCreator,
)

var storageSet = wire.NewSet(
	storage.NewInMemoryCartRepository,
	storage.NewInMemoryItemRepository,
)

var handlerSet = wire.NewSet(
	rest.NewCartHandler,
)

var messagingSet = wire.NewSet(
	getKafkaBrokers,
	messaging.NewKafkaProducer,
)

func InitializeServer(ctx context.Context, cfg Config) (*http.Server, error) {
	wire.Build(
		appSet,
		storageSet,
		messagingSet,
		wire.Bind(new(cart.EventDispatcher), new(*messaging.KafkaEventDispatcher)),
		wire.Bind(new(cart.CartRepository), new(storage.InMemoryCartRepository)),
		wire.Bind(new(cart.ItemRepository), new(storage.InMemoryItemRepository)),
		handlerSet,
		initMetricsExporter,
		metrics.NewOpenCensusRecorder,
		wire.Bind(new(metrics.Recorder), new(metrics.OpenCensusRecorder)),
		getServerAddress,
		rest.NewServer,
	)
	return &http.Server{}, nil
}

func getServerAddress(cfg Config) rest.Address {
	return rest.Address(cfg.ServerPort)
}

func initMetricsExporter(cfg Config) (rest.MetricsHandler, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)

	return exporter, nil
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
