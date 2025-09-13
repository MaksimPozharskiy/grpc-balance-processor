package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicateTx     = errors.New("duplicate transaction")
	ErrNegativeBalance = errors.New("negative balance")
)
