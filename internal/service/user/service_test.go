package user

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"http/internal/domain"
	"http/internal/repository/memory"
	"http/internal/service/user/mocks"
)

func TestService_CreateUser(t *testing.T) {
	type fields struct {
		accountService func() accountService
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr error
	}{
		{
			name: "successfully create user",
			fields: fields{
				accountService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("Create", mock.Anything).Return(nil)
					return accServiceMock
				},
			},
			args: args{
				name: "test",
			},
			want: &domain.User{
				Name: "test",
			},
		},
		{
			name: "invalid name, return failedToCreateUser",
			fields: fields{
				accountService: func() accountService {
					return mocks.NewAccountService(t)
				},
			},
			wantErr: failedToCreateUser,
		},
		{
			name: "fail create account, return failedToCreateAccount",
			fields: fields{
				accountService: func() accountService {
					mockAccountService := mocks.NewAccountService(t)
					mockAccountService.On("Create", mock.Anything).Return(errors.New("fail to create account"))
					return mockAccountService
				},
			},
			args: args{
				name: "test",
			},
			wantErr: failedToCreateAccount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				userRepository: memory.NewUserRepository(),
				accountService: tt.fields.accountService(),
			}
			got, err := service.CreateUser(tt.args.name)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(domain.User{}, "ID")); diff != "" {
				t.Errorf("CreateUser() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	userRepository := memory.NewUserRepository()
	userRepository.Insert(&domain.User{
		ID:   "test",
		Name: "test",
	})

	type fields struct {
		accountsService func() accountService
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "successfully delete user",
			fields: fields{
				accountsService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("DeleteUserAccounts", mock.Anything).Return(nil)
					return accServiceMock
				},
			},
			args: args{
				userID: "test",
			},
		},
		{
			name: "fail to get user to delete, return failedToGetUser",
			fields: fields{
				accountsService: func() accountService {
					return nil
				},
			},
			args: args{
				userID: "invalid",
			},
			wantErr: failedToGetUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				userRepository: userRepository,
				accountService: tt.fields.accountsService(),
			}
			if err := service.DeleteUser(tt.args.userID); !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetAccounts(t *testing.T) {
	userRepository := memory.NewUserRepository()
	userRepository.Insert(&domain.User{
		ID:   "test",
		Name: "test",
	})

	testUserAccounts := []domain.Account{
		{
			ID:      "test",
			UserID:  "test",
			Balance: 2,
		},
	}

	type fields struct {
		accountsService func() accountService
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Account
		wantErr error
	}{
		{
			name: "successfully get accounts",
			fields: fields{
				accountsService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("GetUserAccounts", mock.Anything).Return(testUserAccounts, nil)
					return accServiceMock
				},
			},
			args: args{
				userID: "test",
			},
			want: testUserAccounts,
		},
		{
			name: "invalid userID, return failedToGetUser",
			fields: fields{
				accountsService: func() accountService {
					return nil
				},
			},
			args: args{
				userID: "invalid",
			},
			wantErr: failedToGetUser,
		},
		{
			name: "fail to get accounts, return failedToGetUserAccounts",
			fields: fields{
				accountsService: func() accountService {
					accServiceMock := mocks.NewAccountService(t)
					accServiceMock.On("GetUserAccounts", mock.Anything).Return(nil, errors.New("fail to get accounts"))
					return accServiceMock
				},
			},
			args: args{
				userID: "test",
			},
			wantErr: failedToGetUserAccounts,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				userRepository: userRepository,
				accountService: tt.fields.accountsService(),
			}
			got, err := service.GetAccounts(tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetAccounts() (-want +got):\n%s", diff)
			}
		})
	}
}
