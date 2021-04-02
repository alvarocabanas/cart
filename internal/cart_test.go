package cart_test

import (
	"context"
	"testing"

	cart "github.com/alvarocabanas/cart/internal"

	"github.com/stretchr/testify/require"
)

func TestCart_AddItem(t *testing.T) {
	requireThat := require.New(t)
	ctx := context.Background()

	t.Run("Given a cart and an item", func(t *testing.T) {
		c := cart.New()
		item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
		requireThat.NoError(err)
		t.Run("When we add the item with an amount inferior than one", func(t *testing.T) {
			err = c.AddItem(ctx, item, -2)
			t.Run("Then an error should be returned", func(t *testing.T) {
				requireThat.Equal(cart.ErrWrongAddItemQuantity, err)
			})
		})
		t.Run("When we add the item with an amount superior than one", func(t *testing.T) {
			err = c.AddItem(ctx, item, 1)
			t.Run("Then no error should returned", func(t *testing.T) {
				requireThat.NoError(err)
			})
		})
		t.Run("When we add more items of that kind", func(t *testing.T) {
			err = c.AddItem(ctx, item, 3)
			t.Run("Then no error should returned and the amount should be as expected", func(t *testing.T) {
				requireThat.NoError(err)
				requireThat.Equal(4, c.Lines()[item.UUID()].Quantity())
			})
		})
	})
}

func TestCart_Lines(t *testing.T) {
	requireThat := require.New(t)
	ctx := context.Background()

	t.Run("Given a cart", func(t *testing.T) {
		c := cart.New()
		t.Run("When the cart has 3 lines", func(t *testing.T) {
			item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item, 3)
			requireThat.NoError(err)

			item2, err := cart.NewItem("another-item-id", "book", 400, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item2, 3)
			requireThat.NoError(err)

			item3, err := cart.NewItem("and-another-item-id", "casette", 600, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item3, 3)
			requireThat.NoError(err)

			t.Run("Then the lines returned should be correct", func(t *testing.T) {
				requireThat.Equal(3, len(c.Lines()))
				requireThat.Equal("an-item-id", c.Lines()["an-item-id"].Item().UUID())
			})
		})
	})
}

func TestCart_CalculateTotalPrice(t *testing.T) {
	requireThat := require.New(t)
	ctx := context.Background()

	t.Run("Given a cart with no discounts", func(t *testing.T) {
		c := cart.New()
		t.Run("When the cart has 3 dvds of 100 euros each", func(t *testing.T) {
			item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item, 3)
			requireThat.NoError(err)

			t.Run("Then the price returned should be correct", func(t *testing.T) {
				requireThat.Equal(300, c.CalculateTotalPrice())
			})
		})
		t.Run("When the cart adds 2 books of 50 euros each", func(t *testing.T) {
			item, err := cart.NewItem("another-item-id", "book", 50, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item, 2)
			requireThat.NoError(err)
			t.Run("Then the price returned should be correct", func(t *testing.T) {
				requireThat.Equal(400, c.CalculateTotalPrice())
			})
		})
	})

	t.Run("Given a cart with a 20 EUR discount", func(t *testing.T) {
		discount := func(price int) int {
			price -= 20
			if price == 0 {
				return 0
			}
			return price
		}
		c := cart.New(discount)
		t.Run("When the cart has 3 dvds of 100 euros each", func(t *testing.T) {
			item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
			requireThat.NoError(err)

			err = c.AddItem(ctx, item, 3)
			requireThat.NoError(err)

			t.Run("Then the price returned should be correct", func(t *testing.T) {
				requireThat.Equal(280, c.CalculateTotalPrice())
			})
		})
	})
}
