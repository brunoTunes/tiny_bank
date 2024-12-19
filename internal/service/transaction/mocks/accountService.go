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

// AddBalance provides a mock function with given fields: accountID, balance
func (_m *AccountService) AddBalance(accountID string, balance int) (*domain.Account, error) {
	ret := _m.Called(accountID, balance)

	if len(ret) == 0 {
		panic("no return value specified for AddBalance")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int) (*domain.Account, error)); ok {
		return rf(accountID, balance)
	}
	if rf, ok := ret.Get(0).(func(string, int) *domain.Account); ok {
		r0 = rf(accountID, balance)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(accountID, balance)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: accountID
func (_m *AccountService) Get(accountID string) (*domain.Account, error) {
	ret := _m.Called(accountID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Account, error)); ok {
		return rf(accountID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Account); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountID)
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
