// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	domain "http/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// AccountService is an autogenerated mock type for the accountService type
type AccountService struct {
	mock.Mock
}

// Create provides a mock function with given fields: userID
func (_m *AccountService) Create(userID string) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUserAccounts provides a mock function with given fields: userID
func (_m *AccountService) DeleteUserAccounts(userID string) {
	_m.Called(userID)
}

// GetUserAccounts provides a mock function with given fields: userID
func (_m *AccountService) GetUserAccounts(userID string) ([]domain.Account, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccounts")
	}

	var r0 []domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Account, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Account); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccountService creates a new instance of AccountService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountService {
	mock := &AccountService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
