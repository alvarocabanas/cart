package creator

import (
	"context"

	"go.opencensus.io/trace"

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
	stxt, span := trace.StartSpan(ctx, "cart_creator_add_item")
	defer span.End()

	item, err := c.itemRepository.Get(stxt, dto.ItemID)
	if err != nil {
		return err
	}

	return c.cartRepository.AddItem(stxt, item, dto.Quantity)
}
