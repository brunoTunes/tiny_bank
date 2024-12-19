package domain

import (
	"errors"
	"time"

	"github.com/lithammer/shortuuid/v4"
)

type Transaction struct {
	ID            string
	CreatedAt     time.Time
	FromAccountID string
	ToAccountID   string
	Amount        int
}

func NewTransaction(fromAccountID, toAccountID string, amount int) (*Transaction, error) {
	t := &Transaction{
		ID:            shortuuid.New(),
		CreatedAt:     time.Now(),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}

	return t.validate()
}

var invalidFromAccountError = errors.New("from account ID is required")
var invalidToAccountError = errors.New("to account ID is required")
var invalidAmountError = errors.New("amount must be greater than zero")

func (t *Transaction) validate() (*Transaction, error) {
	if t.FromAccountID == "" {
		return nil, invalidFromAccountError
	}
	if t.ToAccountID == "" {
		return nil, invalidToAccountError
	}
	if t.Amount <= 0 {
		return nil, invalidAmountError
	}

	return t, nil
}
