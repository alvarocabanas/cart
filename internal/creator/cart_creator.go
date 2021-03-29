package creator

import (
	"context"

	cart "github.com/alvarocabanas/cart/internal"
)

type AddItemDTO struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type CartCreator struct {
	cartRepository cart.CartRepository
	itemRepository cart.ItemRepository
}

func NewCartCreator(
	cartRepository cart.CartRepository,
	itemRepository cart.ItemRepository,
) CartCreator {
	return CartCreator{
		cartRepository: cartRepository,
		itemRepository: itemRepository,
	}
}

func (c CartCreator) AddItem(ctx context.Context, dto AddItemDTO) error {
	item, err := c.itemRepository.Get(ctx, dto.ItemID)
	if err != nil {
		return err
	}

	return c.cartRepository.AddItem(ctx, item, dto.Quantity)
}
