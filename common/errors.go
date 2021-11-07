package common

import (
	"errors"
	"fmt"
	terrors "github.com/pkg/errors"
	"os"
)

var (
	ErrorLockAlreadyRequired = errors.New("the lock is occupied")
)

func ErrFmt(err error) {
	fmt.Printf("original error: %T %v\n", terrors.Cause(err), terrors.Cause(err))
	fmt.Printf("stack trace: \n%+v\n", err)
}

func ErrFmtWithExit(err error, code int) {
	ErrFmt(err)
	os.Exit(code)
}
