package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTransaction_validate(t *testing.T) {
	fromAccountID := "1"
	toAccountID := "2"

	type fields struct {
		ID            string
		CreatedAt     time.Time
		FromAccountID *string
		ToAccountID   *string
		Amount        int
		Type          TransactionType
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Transaction
		wantErr error
	}{
		{
			name: "transfer type, valid",
			fields: fields{
				ID:            "1",
				FromAccountID: &fromAccountID,
				ToAccountID:   &toAccountID,
				Amount:        100,
				Type:          Transfer,
			},
			want: &Transaction{
				ID:            "1",
				FromAccountID: &fromAccountID,
				ToAccountID:   &toAccountID,
				Amount:        100,
				Type:          Transfer,
			},
		},
		{
			name: "invalid balance, want invalidBalanceError",
			fields: fields{
				ID:            "1",
				FromAccountID: &fromAccountID,
				ToAccountID:   &toAccountID,
				Amount:        -1,
				Type:          Transfer,
			},
			wantErr: invalidAmountError,
		},
		{
			name: "transfer type, invalid from account, want invalidFromAccountError",
			fields: fields{
				ID:          "1",
				ToAccountID: &toAccountID,
				Amount:      100,
				Type:        Transfer,
			},
			wantErr: invalidFromAccountError,
		},
		{
			name: "transfer type, invalid to account, want invalidToAccountError",
			fields: fields{
				ID:            "1",
				FromAccountID: &fromAccountID,
				Amount:        100,
				Type:          Transfer,
			},
			wantErr: invalidToAccountError,
		},
		{
			name: "deposit type, valid",
			fields: fields{
				ID:          "1",
				ToAccountID: &toAccountID,
				Amount:      100,
				Type:        Deposit,
			},
			want: &Transaction{
				ID:          "1",
				ToAccountID: &toAccountID,
				Amount:      100,
				Type:        Deposit,
			},
		},
		{
			name: "deposit type, invalid to account, want invalidToAccountError",
			fields: fields{
				ID:     "1",
				Amount: 100,
				Type:   Deposit,
			},
			wantErr: invalidToAccountError,
		},
		{
			name: "withdrawal type, invalid from account, want invalidFromAccountError",
			fields: fields{
				ID:     "1",
				Amount: 100,
				Type:   Withdrawal,
			},
			wantErr: invalidFromAccountError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			tr := &Transaction{
				ID:            tt.fields.ID,
				CreatedAt:     tt.fields.CreatedAt,
				FromAccountID: tt.fields.FromAccountID,
				ToAccountID:   tt.fields.ToAccountID,
				Amount:        tt.fields.Amount,
				Type:          tt.fields.Type,
			}
			got, err := tr.validate()
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
