package domain

import (
	"context"

	"github.com/google/uuid"
)

type BalanceRepository interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, error)
	CreateAccount(ctx context.Context, accountID uuid.UUID) error
	ProcessTransaction(ctx context.Context, op *Operation) (*Account, error)
	GetOperationByTxID(ctx context.Context, txID string) (*Operation, error)
}
