package account

import "fmt"

var failedToCreateAccount = fmt.Errorf("failed to create account")
var failedToPersistAccount = fmt.Errorf("failed to persist account")
var failedToGetAccount = fmt.Errorf("failed to get  account")
var failedToAddBalance = fmt.Errorf("failed to add balance")
var invalidAccountID = fmt.Errorf("invalid account ID")
var invalidUserID = fmt.Errorf("invalid user ID")
