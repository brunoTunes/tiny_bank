package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTransaction_validate(t *testing.T) {
	type fields struct {
		ID            string
		CreatedAt     time.Time
		FromAccountID string
		ToAccountID   string
		Amount        int
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Transaction
		wantErr error
	}{
		{
			name: "valid",
			fields: fields{
				ID:            "1",
				FromAccountID: "1",
				ToAccountID:   "1",
				Amount:        100,
			},
			want: &Transaction{
				ID:            "1",
				FromAccountID: "1",
				ToAccountID:   "1",
				Amount:        100,
			},
		},
		{
			name: "invalid from account, want invalidFromAccountError",
			fields: fields{
				ID:          "1",
				ToAccountID: "1",
				Amount:      100,
			},
			wantErr: invalidFromAccountError,
		},
		{
			name: "invalid to account, want invalidToAccountError",
			fields: fields{
				ID:            "1",
				FromAccountID: "1",
				Amount:        100,
			},
			wantErr: invalidToAccountError,
		},
		{
			name: "invalid balance, want invalidBalanceError",
			fields: fields{
				ID:            "1",
				FromAccountID: "1",
				ToAccountID:   "1",
				Amount:        -1,
			},
			wantErr: invalidAmountError,
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
