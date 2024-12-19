package domain

import (
	"errors"
	"time"

	"github.com/lithammer/shortuuid/v4"
)

type Transaction struct {
	ID            string
	CreatedAt     time.Time
	FromAccountID *string
	ToAccountID   *string
	Amount        int
	Type          TransactionType
}

type TransactionType string

const (
	Withdrawal TransactionType = "withdrawal"
	Deposit    TransactionType = "deposit"
	Transfer   TransactionType = "transfer"
)

func (t TransactionType) String() string {
	return string(t)
}

func NewTransfer(fromAccountID, toAccountID string, amount int) (*Transaction, error) {
	t := &Transaction{
		ID:            shortuuid.New(),
		CreatedAt:     time.Now(),
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Amount:        amount,
		Type:          Transfer,
	}

	return t.validate()
}

func NewDeposit(toAccountID string, amount int) (*Transaction, error) {
	t := &Transaction{
		ID:          shortuuid.New(),
		CreatedAt:   time.Now(),
		ToAccountID: &toAccountID,
		Amount:      amount,
		Type:        Deposit,
	}

	return t.validate()
}

func NewWithdrawal(fromAccountID string, amount int) (*Transaction, error) {
	t := &Transaction{
		ID:            shortuuid.New(),
		CreatedAt:     time.Now(),
		FromAccountID: &fromAccountID,
		Amount:        amount,
		Type:          Withdrawal,
	}

	return t.validate()
}

var invalidFromAccountError = errors.New("from account ID is required")
var invalidToAccountError = errors.New("to account ID is required")
var invalidAmountError = errors.New("amount must be greater than zero")
var invalidTransactionType = errors.New("invalid transaction type")

func (t *Transaction) validate() (*Transaction, error) {
	if t.Amount <= 0 {
		return nil, invalidAmountError
	}

	switch {
	case t.Type == Deposit:
		if t.ToAccountID == nil {
			return nil, invalidToAccountError
		}
	case t.Type == Withdrawal:
		if t.FromAccountID == nil {
			return nil, invalidFromAccountError
		}
	case t.Type == Transfer:
		if t.FromAccountID == nil {
			return nil, invalidFromAccountError
		}
		if t.ToAccountID == nil {
			return nil, invalidToAccountError
		}
	default:
		return nil, invalidTransactionType

	}

	return t, nil
}
