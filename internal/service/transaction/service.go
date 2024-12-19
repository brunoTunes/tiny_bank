package transaction

import (
	"errors"
	"sync"
	"time"

	"http/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2 --name=accountService --structname=AccountService --output=mocks/
type accountService interface {
	Get(accountID string) (*domain.Account, error)
	AddBalance(accountID string, balance int) (*domain.Account, error)
}

type transactionRepository interface {
	Insert(transaction *domain.Transaction) (*domain.Transaction, error)
	GetAccountTransactions(accountID string, fromDate time.Time, toDate time.Time) ([]domain.Transaction, error)
}

type Service struct {
	accountService        accountService
	transactionRepository transactionRepository

	// TODO isolate this in it's own package
	mapAccessMutex sync.Mutex
	mutexMap       map[string]*sync.Mutex
}

func NewService(accountService accountService, transactionRepository transactionRepository) *Service {
	return &Service{
		accountService:        accountService,
		transactionRepository: transactionRepository,
		mapAccessMutex:        sync.Mutex{},
		mutexMap:              make(map[string]*sync.Mutex),
	}
}

func (service *Service) Transfer(fromAccountID, toAccountID string, amount int) (*domain.Transaction, error) {
	service.mapAccessMutex.Lock()

	if service.mutexMap[fromAccountID] == nil {
		service.mutexMap[fromAccountID] = &sync.Mutex{}
	}
	service.mutexMap[fromAccountID].Lock()
	defer service.mutexMap[fromAccountID].Unlock()

	if service.mutexMap[toAccountID] == nil {
		service.mutexMap[toAccountID] = &sync.Mutex{}
	}
	service.mutexMap[toAccountID].Lock()
	defer service.mutexMap[toAccountID].Unlock()

	service.mapAccessMutex.Unlock()

	fromAccount, err := service.accountService.AddBalance(fromAccountID, -amount)
	if err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	toAccount, err := service.accountService.AddBalance(toAccountID, amount)
	if err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	transaction, err := domain.NewTransfer(fromAccount.ID, toAccount.ID, amount)
	if err != nil {
		return nil, errors.Join(failedToCreateTransaction, err)
	}

	newTransaction, err := service.transactionRepository.Insert(transaction)
	if err != nil {
		return nil, errors.Join(failedToInsertTransaction, err)
	}

	return newTransaction, nil
}

func (service *Service) Deposit(toAccountID string, amount int) (*domain.Transaction, error) {
	service.mapAccessMutex.Lock()

	if service.mutexMap[toAccountID] == nil {
		service.mutexMap[toAccountID] = &sync.Mutex{}
	}
	service.mutexMap[toAccountID].Lock()
	defer service.mutexMap[toAccountID].Unlock()

	service.mapAccessMutex.Unlock()

	toAccount, err := service.accountService.AddBalance(toAccountID, amount)
	if err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	transaction, err := domain.NewDeposit(toAccount.ID, amount)
	if err != nil {
		return nil, errors.Join(failedToCreateTransaction, err)
	}

	newTransaction, err := service.transactionRepository.Insert(transaction)
	if err != nil {
		return nil, errors.Join(failedToInsertTransaction, err)
	}

	return newTransaction, nil
}

func (service *Service) Withdraw(fromAccountID string, amount int) (*domain.Transaction, error) {
	service.mapAccessMutex.Lock()

	if service.mutexMap[fromAccountID] == nil {
		service.mutexMap[fromAccountID] = &sync.Mutex{}
	}
	service.mutexMap[fromAccountID].Lock()
	defer service.mutexMap[fromAccountID].Unlock()

	service.mapAccessMutex.Unlock()

	fromAccount, err := service.accountService.AddBalance(fromAccountID, -amount)
	if err != nil {
		return nil, errors.Join(failedAddBalance, err)
	}

	transaction, err := domain.NewWithdrawal(fromAccount.ID, amount)
	if err != nil {
		return nil, errors.Join(failedToCreateTransaction, err)
	}

	newTransaction, err := service.transactionRepository.Insert(transaction)
	if err != nil {
		return nil, errors.Join(failedToInsertTransaction, err)
	}

	return newTransaction, nil
}

func (service *Service) GetAccountTransactionHistory(accountID string, fromDate time.Time, toDate time.Time) ([]domain.Transaction, error) {
	if accountID == "" {
		return nil, errors.New("invalid empty account ID")
	}

	transactions, err := service.transactionRepository.GetAccountTransactions(accountID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
