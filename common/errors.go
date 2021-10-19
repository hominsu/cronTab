package common

import "errors"

var (
	ErrorLockAlreadyRequired = errors.New("the lock is occupied")
)
