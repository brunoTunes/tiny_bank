package response

import (
	"time"

	"http/internal/domain"
)

type AccountResponse struct {
	ID        string
	UserID    string
	Balance   int
	DeletedAt *time.Time
}

func AccountResponseFromDomain(account *domain.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		DeletedAt: account.DeletedAt,
	}
}

func AccountTransactionsHistoryFromDomain(transactions []domain.Transaction) []Transaction {
	var transactionsHistory = make([]Transaction, len(transactions))

	for i, transaction := range transactions {
		transactionsHistory[i] = Transaction{
			ID:          transaction.ID,
			FromAccount: transaction.FromAccountID,
			ToAccount:   transaction.ToAccountID,
			Amount:      transaction.Amount,
			CreatedAt:   transaction.CreatedAt,
		}
	}

	return transactionsHistory
}
