package cart

import (
	"context"
	"errors"
	"time"

	"go.opencensus.io/trace"
)

var ErrWrongAddItemQuantity = errors.New("item quantity should be positive on addition")

// Discount is a callback to apply a specific discount to a Cart
// In a possible future, the discount could be moved to the Lines to be able to have different Discounts on each Item
type Discount func(int) int

type Cart struct {
	lines     map[string]*Line
	discounts []Discount
	events    EventList
}

// New Creates a new Cart with a map of Lines
// Each Line has the item type and the quantity of items on that Line
// Discounts is a variadic functional parameter to apply discounts to a Cart price calculation
func New(discounts ...Discount) *Cart {
	return &Cart{
		lines:     make(map[string]*Line),
		discounts: discounts,
	}
}

// Add Item adds items to a New cart Line, RemoveItem is not implemented
func (c *Cart) AddItem(ctx context.Context, item Item, quantity int) error {
	_, span := trace.StartSpan(ctx, "cart_add_item")
	defer span.End()

	if quantity < 1 {
		return ErrWrongAddItemQuantity
	}

	if _, ok := c.lines[item.UUID()]; ok {
		c.lines[item.UUID()].increaseQuantity(quantity)
		return nil
	}

	c.lines[item.UUID()] = NewLine(item, quantity)

	c.events.Record(
		&ItemAddedEvent{
			ItemId:     item.UUID(),
			Quantity:   int32(quantity),
			Type:       "",
			OccurredAt: time.Now().Format(time.RFC3339),
		},
	)
	return nil
}

// Lines returns all the Lines in the cart
func (c *Cart) Lines() map[string]*Line {
	return c.lines
}

// Events returns all events to be dispatched
func (c *Cart) Events() EventList {
	return c.events
}

func (c *Cart) ClearEvents() {
	c.events.Clear()
}

// CalculateTotalPrice calculates the total of the cart, applying discounts to the lines if applicable
func (c Cart) CalculateTotalPrice() int {
	var totalPrice int
	for _, line := range c.lines {
		totalPrice += line.calculatePrice()
	}

	for _, discount := range c.discounts {
		totalPrice = discount(totalPrice)
	}
	return totalPrice
}

type Line struct {
	item     Item
	quantity int
}

// NewLine creates a new Line for an Item and quantity
func NewLine(item Item, quantity int) *Line {
	return &Line{
		item:     item,
		quantity: quantity,
	}
}

// Item returns the Item of the Line
func (l *Line) Item() Item {
	return l.item
}

// Item returns the quantity of Items of the Line
func (l *Line) Quantity() int {
	return l.quantity
}

func (l *Line) increaseQuantity(quantity int) {
	l.quantity += quantity

}

func (l *Line) calculatePrice() int {
	return l.item.price.amount * l.quantity
}
