// +build wireinject

package bootstrap

import (
	"context"
	"net/http"

	cart "github.com/alvarocabanas/cart/internal"
	"github.com/alvarocabanas/cart/internal/creator"
	"github.com/alvarocabanas/cart/internal/getter"
	"github.com/alvarocabanas/cart/internal/io/rest"
	"github.com/alvarocabanas/cart/internal/storage"
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
