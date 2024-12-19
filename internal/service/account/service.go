package account

import (
	"errors"
	"time"

	"http/internal/domain"
)

type accountRepository interface {
	GetUserAccounts(userID string) []domain.Account
	Insert(acc *domain.Account) (*domain.Account, error)
	UpdateBulk(acc []domain.Account)
	Get(accID string) (*domain.Account, error)
}

type Service struct {
	accountRepository accountRepository
}

func NewService(accountRepository accountRepository) *Service {
	return &Service{
		accountRepository: accountRepository,
	}
}

func (service Service) Create(userID string) error {
	acc, err := domain.NewAccount(userID)
	if err != nil {
		return errors.Join(failedToCreateAccount, err)
	}

	acc, err = service.accountRepository.Insert(acc)
	if err != nil {
		return errors.Join(failedToPersistAccount, err)
	}

	return nil
}

func (service Service) DeleteUserAccounts(userID string) {
	accs := service.accountRepository.GetUserAccounts(userID)

	for _, acc := range accs {
		now := time.Now()
		acc.DeletedAt = &now
	}

	service.accountRepository.UpdateBulk(accs)
}

func (service Service) GetUserAccounts(userID string) ([]domain.Account, error) {
	if userID == "" {
		return nil, invalidUserID
	}

	accounts := service.accountRepository.GetUserAccounts(userID)

	return accounts, nil
}

func (service Service) AddBalance(accountID string, balance int) (*domain.Account, error) {
	acc, err := service.accountRepository.Get(accountID)
	if err != nil {
		return nil, errors.Join(failedToGetAccount, err)
	}

	if err := acc.AddBalance(balance); err != nil {
		return nil, errors.Join(failedToAddBalance, err)
	}

	return acc, nil
}

func (service Service) Get(accountID string) (*domain.Account, error) {
	if accountID == "" {
		return nil, invalidAccountID
	}

	return service.accountRepository.Get(accountID)
}
