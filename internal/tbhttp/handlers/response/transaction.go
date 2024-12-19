package response

import (
	"time"

	"http/internal/domain"
)

type Transaction struct {
	ID          string    `json:"id"`
	FromAccount *string   `json:"from_account,omitempty"`
	ToAccount   *string   `json:"to_account,omitempty"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
}

func TransactionFromDomain(transaction *domain.Transaction) Transaction {
	return Transaction{
		ID:          transaction.ID,
		FromAccount: transaction.FromAccountID,
		ToAccount:   transaction.ToAccountID,
		Amount:      transaction.Amount,
		CreatedAt:   transaction.CreatedAt,
		Type:        transaction.Type.String(),
	}
}

func DepositFromDomain(transaction *domain.Transaction) Transaction {
	return Transaction{
		ID:        transaction.ID,
		ToAccount: transaction.ToAccountID,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
		Type:      transaction.Type.String(),
	}
}

func WithdrawFromDomain(transaction *domain.Transaction) Transaction {
	return Transaction{
		ID:          transaction.ID,
		FromAccount: transaction.FromAccountID,
		Amount:      transaction.Amount,
		CreatedAt:   transaction.CreatedAt,
		Type:        transaction.Type.String(),
	}
}

func TransactionsHistoryFromDomain(transactions []domain.Transaction) []Transaction {
	var transactionsHistory = make([]Transaction, len(transactions))

	for i, transaction := range transactions {
		transactionsHistory[i] = Transaction{
			ID:          transaction.ID,
			FromAccount: transaction.FromAccountID,
			ToAccount:   transaction.ToAccountID,
			Amount:      transaction.Amount,
			CreatedAt:   transaction.CreatedAt,
			Type:        transaction.Type.String(),
		}
	}

	return transactionsHistory
}
