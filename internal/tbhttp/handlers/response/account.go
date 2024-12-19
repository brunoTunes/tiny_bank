package response

import (
	"time"

	"http/internal/domain"
)

type AccountResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Balance   int        `json:"balance"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func AccountResponseFromDomain(account *domain.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		DeletedAt: account.DeletedAt,
	}
}
