package cart

import "errors"

var (
	ErrEmptyItemID   = errors.New("item id cannot be empty")
	ErrEmptyItemName = errors.New("item name cannot be empty")
	ErrItemNotFound  = errors.New("item not found")
)

type Item struct {
	id    ItemUUID
	name  string
	price Money
}

// NewItem creates item entities, the creation constraints are checked in this factory method
// or the ones in its Value Objects, Money and ItemUUID
func NewItem(id, name string, price int, currency string) (Item, error) {
	itemID, err := NewItemUUID(id)
	if err != nil {
		return Item{}, err
	}

	if name == "" {
		return Item{}, ErrEmptyItemName
	}

	itemPrice, err := NewMoney(price, currency)
	if err != nil {
		return Item{}, err
	}

	return Item{
		id:    itemID,
		name:  name,
		price: itemPrice,
	}, nil
}

type ItemUUID string

func NewItemUUID(id string) (ItemUUID, error) {
	if id == "" {
		return "", ErrEmptyItemID
	}

	return ItemUUID(id), nil
}

func (i Item) UUID() string {
	return string(i.id)
}
