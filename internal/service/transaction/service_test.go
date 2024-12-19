package transaction

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lithammer/shortuuid/v4"
	"http/internal/domain"
	"http/internal/repository/memory"
	"http/internal/service/transaction/mocks"
)

func TestService_Transfer(t *testing.T) {
	transactionRepository := memory.NewTransactionRepository()
	fromAccount := &domain.Account{
		ID:      "1",
		UserID:  "1",
		Balance: 100,
	}

	toAccount := &domain.Account{
		ID:      "2",
		UserID:  "2",
		Balance: 0,
	}

	type fields struct {
		accountService func() accountService
	}
	type args struct {
		fromAccountID string
		toAccountID   string
		amount        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Transaction
		wantErr error
	}{
		{
			name: "successful transfer",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Get", fromAccount.ID).Return(fromAccount, nil)
					accServiceMock.On("Get", toAccount.ID).Return(toAccount, nil)
					return accServiceMock
				},
			},
			args: args{
				fromAccountID: fromAccount.ID,
				toAccountID:   toAccount.ID,
				amount:        100,
			},
			want: &domain.Transaction{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        100,
			},
		},
		{
			name: "failed transfer, fromAccount without enough balance, return failedAddBalance",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Get", fromAccount.ID).Return(fromAccount, nil)
					return accServiceMock
				},
			},
			args: args{
				fromAccountID: fromAccount.ID,
				toAccountID:   toAccount.ID,
				amount:        1000,
			},
			wantErr: failedAddBalance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountService:        tt.fields.accountService(),
				transactionRepository: transactionRepository,
			}
			got, err := service.Transfer(tt.args.fromAccountID, tt.args.toAccountID, tt.args.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(domain.Transaction{}, "ID", "CreatedAt")); diff != "" {
				t.Errorf("Transfer() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestService_GetAccountTransactionHistory(t *testing.T) {
	transactionRepository := memory.NewTransactionRepository()
	now := time.Now()
	fourDaysAgo := now.Add(-24 * time.Hour * 4)
	fiveDaysAgo := now.Add(-24 * time.Hour * 5)

	t1, _ := transactionRepository.Insert(&domain.Transaction{
		ID:            shortuuid.New(),
		FromAccountID: "1",
		ToAccountID:   "2",
		Amount:        100,
		CreatedAt:     now,
	})
	t2, _ := transactionRepository.Insert(&domain.Transaction{
		ID:            shortuuid.New(),
		FromAccountID: "2",
		ToAccountID:   "1",
		Amount:        100,
		CreatedAt:     now.AddDate(0, 0, -1),
	})
	t3, _ := transactionRepository.Insert(&domain.Transaction{
		ID:            shortuuid.New(),
		FromAccountID: "1",
		ToAccountID:   "2",
		Amount:        100,
		CreatedAt:     now.AddDate(0, 0, -2),
	})
	t4, _ := transactionRepository.Insert(&domain.Transaction{
		ID:            shortuuid.New(),
		FromAccountID: "2",
		ToAccountID:   "1",
		Amount:        100,
		CreatedAt:     now.AddDate(0, 0, -3),
	})

	type fields struct {
		accountService func() accountService
	}
	type args struct {
		accountID string
		fromDate  time.Time
		toDate    time.Time
		limit     int
		offset    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Transaction
		wantErr error
	}{
		{
			name: "transaction history with today's transactions",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Get", "1").Return(&domain.Account{ID: "1"}, nil)
					return accServiceMock
				},
			},
			args: args{
				accountID: "1",
				fromDate:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
				toDate:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC),
				limit:     10,
				offset:    0,
			},
			want: []domain.Transaction{
				*t1,
			},
		},
		{
			name: "transaction history with last 4 days' transactions",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Get", "1").Return(&domain.Account{ID: "1"}, nil)
					return accServiceMock
				},
			},
			args: args{
				accountID: "1",
				fromDate:  time.Date(fourDaysAgo.Year(), fourDaysAgo.Month(), fourDaysAgo.Day(), 0, 0, 0, 0, time.UTC),
				toDate:    time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC),
			},
			want: []domain.Transaction{
				*t4, *t3, *t2, *t1,
			},
		},
		{
			name: "transaction history from 5 days ago to 4 days ago, empty list",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Get", "1").Return(&domain.Account{ID: "1"}, nil)
					return accServiceMock
				},
			},
			args: args{
				accountID: "1",
				fromDate:  time.Date(fiveDaysAgo.Year(), fiveDaysAgo.Month(), fiveDaysAgo.Day(), 0, 0, 0, 0, time.UTC),
				toDate:    time.Date(fourDaysAgo.Year(), fourDaysAgo.Month(), fourDaysAgo.Day(), 23, 59, 59, 0, time.UTC),
			},
			want: []domain.Transaction{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				accountService:        tt.fields.accountService(),
				transactionRepository: transactionRepository,
			}
			got, err := service.GetAccountTransactionHistory(tt.args.accountID, tt.args.fromDate, tt.args.toDate)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAccountTransactionHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetAccountTransactionHistory() (-want +got):\n%s", diff)
			}
		})
	}
}
