// +build wireinject

package bootstrap

import (
	"context"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	cart "github.com/alvarocabanas/cart/internal"
	"github.com/alvarocabanas/cart/internal/creator"
	"github.com/alvarocabanas/cart/internal/getter"
	"github.com/alvarocabanas/cart/internal/io/rest"
	"github.com/alvarocabanas/cart/internal/observability"
	"github.com/alvarocabanas/cart/internal/storage"
	"github.com/google/wire"
	"go.opencensus.io/stats/view"
)

// In future iterations the config will come from Environment variables
type Config struct {
	ServerPort string `mapstructure:"server_port"`
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

func InitializeServer(ctx context.Context, cfg Config) (*http.Server, error) {
	wire.Build(
		appSet,
		storageSet,
		wire.Bind(new(cart.CartRepository), new(storage.InMemoryCartRepository)),
		wire.Bind(new(cart.ItemRepository), new(storage.InMemoryItemRepository)),
		handlerSet,
		initializeMetricsExporter,
		observability.NewOpenCensusMetricsTracker,
		wire.Bind(new(observability.MetricsTracker), new(observability.OpenCensusMetricsTracker)),
		getServerAddress,
		rest.NewServer,
	)
	return &http.Server{}, nil
}

func getServerAddress(cfg Config) rest.Address {
	return rest.Address(cfg.ServerPort)
}

func initializeMetricsExporter(cfg Config) (rest.MetricsHandler, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)

	return exporter, nil
}
