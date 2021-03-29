package cart

import "context"

// This application only allows one cart, in future iterations, the cart should be first Created,
// and the UUID of the cart should be passed to the other calls
type CartRepository interface {
	AddItem(ctx context.Context, item Item, quantity int) error
	Get(ctx context.Context) Cart
}

// This application only has Get Items but in future iterations there would also be a Save method to add new ones
type ItemRepository interface {
	Get(ctx context.Context, itemID string) (Item, error)
}
