package storage

import (
	"context"

	cart "github.com/alvarocabanas/cart/internal"
)

type InMemoryItemRepository struct {
	items map[string]cart.Item
}

// Mocked data
func NewInMemoryItemRepository() InMemoryItemRepository {
	itemA, _ := cart.NewItem("dvd", "dvd", 100, "EUR")
	itemB, _ := cart.NewItem("book", "book", 60, "EUR")
	itemC, _ := cart.NewItem("casette", "casette", 40, "EUR")
	mockItems := map[string]cart.Item{
		itemA.UUID(): itemA,
		itemB.UUID(): itemB,
		itemC.UUID(): itemC,
	}
	return InMemoryItemRepository{
		items: mockItems,
	}
}

func (r InMemoryItemRepository) Get(_ context.Context, itemID string) (cart.Item, error) {
	if item, ok := r.items[itemID]; ok {
		return item, nil
	}
	return cart.Item{}, cart.ErrItemNotFound
}
