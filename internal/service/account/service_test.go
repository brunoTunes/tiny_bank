package account

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"http/internal/domain"
	"http/internal/repository/memory"
)

func TestService_AddBalance(t *testing.T) {
	accountRepository := memory.NewAccountRepository()
	accountRepository.Insert(&domain.Account{
		ID:      "1",
		Balance: 0,
	})
	accountRepository.Insert(&domain.Account{
		ID:      "withBalance",
		Balance: 100,
	})

	type args struct {
		accountID string
		balance   int
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Account
		wantErr error
	}{
		{
			name: "successfully add balance",
			args: args{
				accountID: "1",
				balance:   100,
			},
			want: &domain.Account{
				ID:      "1",
				Balance: 100,
			},
		},
		{
			name: "successfully subtract balance",
			args: args{
				accountID: "withBalance",
				balance:   -100,
			},
			want: &domain.Account{
				ID:      "withBalance",
				Balance: 0,
			},
		},
		{
			name: "invalid account id, return failedToGetAccount",
			args: args{
				accountID: "invalid",
				balance:   100,
			},
			wantErr: failedToGetAccount,
		},
		{
			name: "invalid balance, return failedToAddBalance",
			args: args{
				accountID: "1",
				balance:   -1000,
			},
			wantErr: failedToAddBalance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountRepository: accountRepository,
			}
			got, err := service.AddBalance(tt.args.accountID, tt.args.balance)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AddBalance() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	accountRepository := memory.NewAccountRepository()

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "successfully create account",
			args: args{
				userID: "1",
			},
		},
		{
			name: "invalid user id, return failedToCreateAccount",
			args: args{
				userID: "",
			},
			wantErr: failedToCreateAccount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountRepository: accountRepository,
			}
			if err := service.Create(tt.args.userID); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	accountRepository := memory.NewAccountRepository()
	accountRepository.Insert(&domain.Account{
		ID:      "1",
		Balance: 0,
	})

	type args struct {
		accountID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Account
		wantErr error
	}{
		{
			name: "successfully get account",
			args: args{
				accountID: "1",
			},
			want: &domain.Account{
				ID:      "1",
				Balance: 0,
			},
		},
		{
			name: "empty account id, return invalidAccountID",
			args: args{
				accountID: "",
			},
			wantErr: invalidAccountID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountRepository: accountRepository,
			}
			got, err := service.Get(tt.args.accountID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Get() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestService_GetUserAccounts(t *testing.T) {
	userAccount := domain.Account{
		ID:      "1",
		UserID:  "1",
		Balance: 0,
	}
	accountRepository := memory.NewAccountRepository()
	accountRepository.Insert(&userAccount)

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.Account
		wantErr error
	}{
		{
			name: "successfully get accounts",
			args: args{
				userID: "1",
			},
			want: []domain.Account{
				userAccount,
			},
		},
		{
			name: "empty user id, return invalidUserID",
			args: args{
				userID: "",
			},
			wantErr: invalidUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountRepository: accountRepository,
			}
			got, err := service.GetUserAccounts(tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetUserAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserAccounts() (-want +got):\n%s", diff)
			}
		})
	}
}
