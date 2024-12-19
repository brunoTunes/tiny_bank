package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"http/internal/domain"
)

type TransactionRepository struct {
	transactions          map[string]*domain.Transaction
	transactionsDateIndex map[time.Time][]string
	mutex                 sync.RWMutex
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions:          make(map[string]*domain.Transaction),
		transactionsDateIndex: make(map[time.Time][]string),
		mutex:                 sync.RWMutex{},
	}
}

func (repo *TransactionRepository) Insert(transaction *domain.Transaction) (*domain.Transaction, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if _, ok := repo.transactions[transaction.ID]; ok {
		return nil, errors.New("transaction with id already exists")
	}

	repo.transactions[transaction.ID] = transaction
	repo.transactionsDateIndex[transaction.CreatedAt] = append(repo.transactionsDateIndex[transaction.CreatedAt], transaction.ID)

	return repo.transactions[transaction.ID], nil
}

func (repo *TransactionRepository) GetAccountTransactions(accountID string, fromDate time.Time, toDate time.Time) ([]domain.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var dateKeys []time.Time

	for key, _ := range repo.transactionsDateIndex {
		dateKeys = append(dateKeys, key)
	}

	sort.Slice(dateKeys, func(i, j int) bool {
		return dateKeys[i].Before(dateKeys[j])
	})

	var transactions = make([]domain.Transaction, 0)
	for _, date := range dateKeys {
		if date.Before(fromDate) {
			continue
		}
		if date.After(toDate) {
			break
		}

		for _, transactionID := range repo.transactionsDateIndex[date] {
			if repo.transactions[transactionID].ToAccountID == accountID || repo.transactions[transactionID].FromAccountID == accountID {
				transactions = append(transactions, *repo.transactions[transactionID])
			}
		}
	}

	return transactions, nil
}
