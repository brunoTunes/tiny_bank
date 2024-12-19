package user

import "errors"

var failedToCreateUser = errors.New("failed to create user")
var failedToCreateAccount = errors.New("failed to create account")
var failedToPersistUser = errors.New("failed to persist user")
var failedToGetUser = errors.New("failed to get user")
var failedToUpdateUser = errors.New("failed to update user")
var failedToGetUserAccounts = errors.New("failed to get user accounts")
