package getter

import (
	"context"

	cart "github.com/alvarocabanas/cart/internal"
)

// This in future iterations could go in a transformer to isolate the responsability
type GetCartStatusDTO struct {
	Items []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
	TotalPrice int `json:"total_price"`
}

type CartGetter struct {
	cartRepository cart.CartRepository
}

func NewCartGetter(
	cartRepository cart.CartRepository,
) CartGetter {
	return CartGetter{
		cartRepository: cartRepository,
	}
}

func (g CartGetter) Get(ctx context.Context) GetCartStatusDTO {
	var dto GetCartStatusDTO

	c := g.cartRepository.Get(ctx)
	for _, line := range c.Lines() {
		dto.Items = append(dto.Items, struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}{Id: line.Item().UUID(), Quantity: line.Quantity()})
	}

	dto.TotalPrice = c.CalculateTotalPrice()
	return dto
}
