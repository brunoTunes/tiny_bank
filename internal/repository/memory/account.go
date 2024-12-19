package memory

import (
	"errors"
	"sync"

	"http/internal/domain"
)

type AccountRepository struct {
	transactions map[string]*domain.Transaction
	accounts     map[string]*domain.Account
	mutex        sync.RWMutex
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: make(map[string]*domain.Account),
		mutex:    sync.RWMutex{},
	}
}

func (repo *AccountRepository) Insert(account *domain.Account) (*domain.Account, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if repo.accounts[account.ID] != nil {
		return nil, errors.New("account with id already exists")
	}

	repo.accounts[account.ID] = account

	return repo.accounts[account.ID], nil
}

func (repo *AccountRepository) Get(accID string) (*domain.Account, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	account, ok := repo.accounts[accID]

	if !ok {
		return nil, errors.New("account with id does not exists")
	}

	return account, nil
}

func (repo *AccountRepository) GetUserAccounts(userID string) []domain.Account {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var accounts []domain.Account

	for _, account := range repo.accounts {
		if account.UserID == userID {
			accounts = append(accounts, *account)
		}
	}

	return accounts
}

func (repo *AccountRepository) UpdateBulk(accounts []domain.Account) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, account := range accounts {
		repo.accounts[account.ID] = &account
	}
}
