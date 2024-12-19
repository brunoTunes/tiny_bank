package response

import (
	"time"

	"http/internal/domain"
)

type UserResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func UserResponseFromDomain(domainUser *domain.User) UserResponse {
	return UserResponse{
		ID:        domainUser.ID,
		Name:      domainUser.Name,
		DeletedAt: domainUser.DeletedAt,
	}
}

func AccountsResponseFromDomain(domainAccounts []domain.Account) []AccountResponse {
	accounts := make([]AccountResponse, len(domainAccounts))

	for i, domainAccount := range domainAccounts {
		accounts[i] = AccountResponse{
			ID:        domainAccount.ID,
			UserID:    domainAccount.UserID,
			Balance:   domainAccount.Balance,
			DeletedAt: domainAccount.DeletedAt,
		}
	}

	return accounts
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

func GetUsersFromDomain(users []domain.User) ListUsersResponse {
	var listUsers = make([]UserResponse, len(users))

	for i, user := range users {
		listUsers[i] = UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			DeletedAt: user.DeletedAt,
		}
	}

	return ListUsersResponse{Users: listUsers}
}
