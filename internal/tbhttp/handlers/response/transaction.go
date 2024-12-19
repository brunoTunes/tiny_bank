package response

import (
	"time"

	"http/internal/domain"
)

type Transaction struct {
	ID          string    `json:"id"`
	FromAccount string    `json:"from_account"`
	ToAccount   string    `json:"to_account"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func TransactionFromDomain(transaction *domain.Transaction) Transaction {
	return Transaction{
		ID:          transaction.ID,
		FromAccount: transaction.FromAccountID,
		ToAccount:   transaction.ToAccountID,
		Amount:      transaction.Amount,
		CreatedAt:   transaction.CreatedAt,
	}
}
