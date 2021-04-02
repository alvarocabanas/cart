package storage

import (
	"context"
	"testing"

	cart "github.com/alvarocabanas/cart/internal"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCartRepository_AddItem(t *testing.T) {
	t.Run("Given a repository", func(t *testing.T) {
		repository := NewInMemoryCartRepository()
		item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
		require.NoError(t, err)
		t.Run("When an item is added", func(t *testing.T) {
			err := repository.UpdateLine(context.Background(), cart.NewLine(item, 5))
			t.Run("Then everything works as expected with no race conditions", func(t *testing.T) {
				require.NoError(t, err)
			})
		})
	})
}
