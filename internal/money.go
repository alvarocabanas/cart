package cart

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidAmount   = errors.New("amount cannot be negative")
)

type Money struct {
	amount   int
	currency Currency
}

// NewMoney creates Money value objects, the Currency must be a valid 3 letter ISO Code and the amount be positive
func NewMoney(amount int, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrInvalidAmount
	}

	moneyCurrency, err := NewCurrency(currency)
	if err != nil {
		return Money{}, err
	}

	return Money{
		amount:   amount,
		currency: moneyCurrency,
	}, nil
}

type Currency string

// Regex checking if the isoCode is 3 letters
var isoCodeRegex = regexp.MustCompile("[A-Za-z]{3}")

func NewCurrency(isoCode string) (Currency, error) {
	if !isoCodeRegex.MatchString(isoCode) {
		return "", ErrInvalidCurrency
	}

	return Currency(isoCode), nil
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currency() string {
	return string(m.currency)
}
