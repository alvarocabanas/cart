package storage

import (
	"context"

	"go.opencensus.io/trace"

	cart "github.com/alvarocabanas/cart/internal"
)

// In this repository there should be a map of Carts, but to simplify it for this example, I only add one cart
type InMemoryCartRepository struct {
	cart cart.Cart
}

func NewInMemoryCartRepository() InMemoryCartRepository {
	return InMemoryCartRepository{
		cart: cart.New(),
	}
}

func (r InMemoryCartRepository) AddItem(ctx context.Context, item cart.Item, quantity int) error {
	_, span := trace.StartSpan(ctx, "repository_add_item")
	defer span.End()
	return r.cart.AddItem(item, quantity)
}

func (r InMemoryCartRepository) Get(_ context.Context) cart.Cart {
	return r.cart
}
