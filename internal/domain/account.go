package domain

import (
	"errors"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"http/internal/tberrors"
)

type Account struct {
	ID        string
	UserID    string
	Balance   int
	DeletedAt *time.Time
}

func NewAccount(userID string) (*Account, error) {
	acc := Account{
		ID:      shortuuid.New(),
		UserID:  userID,
		Balance: 0,
	}

	return acc.validate()
}

var emptyAccountUserIDError = tberrors.NewValidationError("invalid empty account user id", "UserID")

func (acc *Account) validate() (*Account, error) {
	if acc.UserID == "" {
		return nil, emptyAccountUserIDError
	}

	return acc, nil
}

var negativeBalanceError = errors.New("adding balance results in negative balance")

func (acc *Account) AddBalance(balance int) error {
	if acc.Balance+balance < 0 {
		return negativeBalanceError
	}

	acc.Balance += balance

	return nil
}
