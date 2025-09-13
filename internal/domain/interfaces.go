package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BalanceRepository interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, error)
	CreateAccount(ctx context.Context, accountID uuid.UUID) error
	ProcessTransaction(ctx context.Context, op *Operation) (*Account, error)
	GetOperationByTxID(ctx context.Context, txID string) (*Operation, error)
}

type BalanceService interface {
	Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error)
	GetBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceResponse, error)
}

type ProcessRequest struct {
	AccountID uuid.UUID
	Source    Source
	State     State
	Amount    decimal.Decimal
	TxID      string
}

type ProcessResponse struct {
	TxID      string
	Status    ProcessStatus
	Balance   decimal.Decimal
	Timestamp time.Time
}

type GetBalanceRequest struct {
	AccountID uuid.UUID
}

type GetBalanceResponse struct {
	Balance   decimal.Decimal
	UpdatedAt time.Time
}

type ProcessStatus int

const (
	StatusOK ProcessStatus = iota
	StatusAlreadyProcessed
	StatusRejectedNegative
)
