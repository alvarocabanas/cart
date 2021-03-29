package cart_test

import (
	cart "github.com/alvarocabanas/cart/internal"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	requireThat := require.New(t)

	t.Run("Given a new Item", func(t *testing.T) {
		t.Run("When the id is empty", func(t *testing.T) {
			item, err := cart.NewItem("", "dvd", 100, "EUR")
			t.Run("Then the item should not be created and return error", func(t *testing.T) {
				requireThat.Empty(item)
				requireThat.Equal(cart.ErrEmptyItemID, err)
			})
		})
		t.Run("When the name is empty", func(t *testing.T) {
			item, err := cart.NewItem("new-id", "", 100, "EUR")
			t.Run("Then the item should not be created and return error", func(t *testing.T) {
				requireThat.Empty(item)
				requireThat.Equal(cart.ErrEmptyItemName, err)
			})
		})
		t.Run("When the price amount is negative", func(t *testing.T) {
			item, err := cart.NewItem("new-id", "dvd", -1, "EUR")
			t.Run("Then the item should not be created and return error", func(t *testing.T) {
				requireThat.Empty(item)
				requireThat.Equal(cart.ErrInvalidAmount, err)
			})
		})
		t.Run("When the price currency is not a valid ISO code", func(t *testing.T) {
			item, err := cart.NewItem("new-id", "dvd", 10, "M33")
			t.Run("Then the item should not be created and return error", func(t *testing.T) {
				requireThat.Empty(item)
				requireThat.Equal(cart.ErrInvalidCurrency, err)
			})
		})
		t.Run("When everything is correct", func(t *testing.T) {
			item, err := cart.NewItem("new-id", "dvd", 10, "EUR")
			t.Run("Then the item should be created and no error returned", func(t *testing.T) {
				requireThat.NotEmpty(item)
				requireThat.NoError(err)
			})
		})
	})
}
