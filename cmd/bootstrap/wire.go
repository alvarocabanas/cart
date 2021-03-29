// +build wireinject

package bootstrap

import (
	cart "cart/internal"
	"cart/internal/creator"
	"cart/internal/getter"
	"cart/internal/io/rest"
	"cart/internal/storage"
	"context"
	"net/http"

	"github.com/google/wire"
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
		getServerAddress,
		rest.NewServer,
	)
	return &http.Server{}, nil
}

func getServerAddress(cfg Config) rest.Address {
	return rest.Address(cfg.ServerPort)
}
