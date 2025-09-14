package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Source string

const (
	SourceGame    Source = "game"
	SourcePayment Source = "payment"
	SourceService Source = "service"
)

type State string

const (
	StateDeposit  State = "deposit"
	StateWithdraw State = "withdraw"
)

type Account struct {
	ID        uuid.UUID
	Balance   decimal.Decimal
	UpdatedAt time.Time
}

type Operation struct {
	ID         int64
	TxID       string
	AccountID  uuid.UUID
	Source     Source
	State      State
	Amount     decimal.Decimal
	CreatedAt  time.Time
	Applied    bool
	CanceledAt *time.Time
	CancelNote *string
}
