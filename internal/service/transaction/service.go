package transaction

import (
	"errors"
	"time"

	"http/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2 --name=accountService --structname=AccountService --output=mocks/
type accountService interface {
	Get(accountID string) (*domain.Account, error)
}

type transactionRepository interface {
	Insert(transaction *domain.Transaction) (*domain.Transaction, error)
	GetAccountTransactions(accountID string, fromDate time.Time, toDate time.Time) ([]domain.Transaction, error)
}

type Service struct {
	accountService        accountService
	transactionRepository transactionRepository
}

func NewService(accountService accountService, transactionRepository transactionRepository) *Service {
	return &Service{
		accountService:        accountService,
		transactionRepository: transactionRepository,
	}
}

func (service Service) Transfer(fromAccountID, toAccountID string, amount int) (*domain.Transaction, error) {
	transaction, err := domain.NewTransaction(fromAccountID, toAccountID, amount)
	if err != nil {
		return nil, errors.Join(failedToCreateTransaction, err)
	}

	fromAccount, err := service.accountService.Get(transaction.FromAccountID)
	if err != nil {
		return nil, errors.Join(failedToGetAccount, err)
	}

	if err := fromAccount.AddBalance(-amount); err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	toAccount, err := service.accountService.Get(transaction.ToAccountID)
	if err != nil {
		return nil, errors.Join(failedToGetAccount, err)
	}

	if err := toAccount.AddBalance(amount); err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	newTransaction, err := service.transactionRepository.Insert(transaction)
	if err != nil {
		return nil, errors.Join(failedToInsertTransaction, err)
	}

	return newTransaction, nil
}

func (service Service) GetAccountTransactionHistory(accountID string, fromDate time.Time, toDate time.Time) ([]domain.Transaction, error) {
	if accountID == "" {
		return nil, errors.New("invalid empty account ID")
	}

	acc, err := service.accountService.Get(accountID)
	if err != nil {
		return nil, err
	}

	transactions, err := service.transactionRepository.GetAccountTransactions(acc.ID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
