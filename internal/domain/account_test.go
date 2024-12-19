package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAccount_validate(t *testing.T) {
	type fields struct {
		ID        string
		UserID    string
		Balance   int
		DeletedAt *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Account
		wantErr error
	}{
		{
			name: "valid account",
			fields: fields{
				ID:      "1",
				UserID:  "1",
				Balance: 0,
			},
			want: &Account{
				ID:      "1",
				UserID:  "1",
				Balance: 0,
			},
		},
		{
			name: "invalid account ownerid, want emptyAccountUserIDError",
			fields: fields{
				ID:      "1",
				UserID:  "",
				Balance: 0,
			},
			wantErr: emptyAccountUserIDError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &Account{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Balance:   tt.fields.Balance,
				DeletedAt: tt.fields.DeletedAt,
			}

			got, err := acc.validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("validate() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAccount_AddBalance(t *testing.T) {
	type fields struct {
		ID        string
		UserID    string
		Balance   int
		DeletedAt *time.Time
	}
	type args struct {
		balance int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantBalance int
		wantErr     error
	}{
		{
			name: "add balance",
			fields: fields{
				ID:      "1",
				UserID:  "1",
				Balance: 0,
			},
			args: args{
				balance: 1,
			},
			wantBalance: 1,
		},
		{
			name: "subtract balance",
			fields: fields{
				ID:      "1",
				UserID:  "1",
				Balance: 1,
			},
			args: args{
				balance: -1,
			},
			wantBalance: 0,
		},
		{
			name: "subtract balance from account without sufficient balance",
			fields: fields{
				ID:      "1",
				UserID:  "1",
				Balance: 0,
			},
			args: args{
				balance: -1,
			},
			wantErr: negativeBalanceError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &Account{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Balance:   tt.fields.Balance,
				DeletedAt: tt.fields.DeletedAt,
			}
			if err := acc.AddBalance(tt.args.balance); !errors.Is(err, tt.wantErr) {
				t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
			}

			if acc.Balance != tt.wantBalance {
				t.Errorf("AddBalance() got = %v, want %v", acc.Balance, tt.wantBalance)
			}
		})
	}
}
