package creator

import (
	"context"

	"go.opencensus.io/trace"

	cartpkg "github.com/alvarocabanas/cart/internal"
)

type AddItemDTO struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type CartCreator struct {
	cartRepository  cartpkg.CartRepository
	itemRepository  cartpkg.ItemRepository
	eventDispatcher cartpkg.EventDispatcher
}

func NewCartCreator(
	cartRepository cartpkg.CartRepository,
	itemRepository cartpkg.ItemRepository,
	eventDispatcher cartpkg.EventDispatcher,
) CartCreator {
	return CartCreator{
		cartRepository:  cartRepository,
		itemRepository:  itemRepository,
		eventDispatcher: eventDispatcher,
	}
}

func (c CartCreator) AddItem(ctx context.Context, dto AddItemDTO) error {
	stxt, span := trace.StartSpan(ctx, "cart_creator_add_item")
	defer span.End()

	item, err := c.itemRepository.Get(stxt, dto.ItemID)
	if err != nil {
		return err
	}

	cart := c.cartRepository.Get(stxt)
	err = cart.AddItem(ctx, item, dto.Quantity)
	if err != nil {
		return err
	}

	err = c.cartRepository.UpdateLine(stxt, cart.Lines()[item.UUID()])
	if err != nil {
		return err
	}

	for _, event := range cart.Events() {
		err = c.eventDispatcher.Dispatch(stxt, cartpkg.EventsTopic, item.UUID(), event)
		if err != nil {
			return err
		}
	}
	cart.ClearEvents()
	return nil
}
