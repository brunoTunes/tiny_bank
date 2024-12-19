package transaction

import (
	"errors"
)

var failedToCreateTransaction = errors.New("failed to create transaction")
var failedToGetAccount error = errors.New("failed to get account")
var failedAddBalance error = errors.New("failed to add balance")
var failedToInsertTransaction = errors.New("failed to insert transaction")
