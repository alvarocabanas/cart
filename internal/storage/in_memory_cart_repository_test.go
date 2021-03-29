package storage

import (
	cart "cart/internal"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryCartRepository_AddItem(t *testing.T) {
	t.Run("Given a repository", func(t *testing.T) {
		repository := NewInMemoryCartRepository()
		item, err := cart.NewItem("an-item-id", "dvd", 100, "EUR")
		require.NoError(t, err)
		t.Run("When an item is added", func(t *testing.T) {
			err := repository.AddItem(context.Background(), item, 5)
			t.Run("Then everything works as expected with no race conditions", func(t *testing.T) {
				require.NoError(t, err)
			})
		})
	})
}
