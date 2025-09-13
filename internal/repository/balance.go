package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/db"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	sqlCreateAccount = `
INSERT INTO accounts (id, balance) VALUES ($1, 0)
ON CONFLICT (id) DO NOTHING
`

	sqlInsertOperation = `
INSERT INTO operations (tx_id, account_id, source, state, amount)
VALUES ($1, $2, $3::source_t, $4::state_t, $5::numeric)
ON CONFLICT (tx_id) DO NOTHING
`

	sqlUpdateBalance = `
UPDATE accounts
   SET balance = balance + $1::numeric,
       updated_at = now()
 WHERE id = $2
   AND balance + $1::numeric >= 0
RETURNING balance, updated_at
`

	sqlSelectBalance = `
SELECT balance, updated_at FROM accounts WHERE id = $1
`
)

type BalanceRepository struct {
	db *db.DB
}

func NewBalanceRepository(database *db.DB) *BalanceRepository {
	return &BalanceRepository{db: database}
}

func (r *BalanceRepository) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	query := `SELECT id, balance, updated_at FROM accounts WHERE id = $1`

	var acc domain.Account
	var balanceStr string
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&acc.ID, &balanceStr, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get account: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get account query: %w", err)
	}

	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		return nil, fmt.Errorf("parse balance: %w", err)
	}
	acc.Balance = balance

	return &acc, nil
}

func (r *BalanceRepository) CreateAccount(ctx context.Context, accountID uuid.UUID) error {
	query := `INSERT INTO accounts (id, balance) VALUES ($1, 0) ON CONFLICT (id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	return nil
}

func (r *BalanceRepository) GetOperationByTxID(ctx context.Context, txID string) (*domain.Operation, error) {
	query := `
		SELECT id, tx_id, account_id, source, state, amount, created_at, applied, canceled_at, cancel_note
		FROM operations
		WHERE tx_id = $1
	`

	var op domain.Operation
	err := r.db.QueryRowContext(ctx, query, txID).Scan(
		&op.ID, &op.TxID, &op.AccountID, &op.Source, &op.State, &op.Amount,
		&op.CreatedAt, &op.Applied, &op.CanceledAt, &op.CancelNote,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get operation by tx_id: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get operation query: %w", err)
	}

	return &op, nil
}

func selectAccount(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Account, error) {
	var s string
	var t time.Time
	if err := tx.QueryRowContext(ctx, sqlSelectBalance, id).Scan(&s, &t); err != nil {
		return nil, err
	}
	bal, err := decimal.NewFromString(s)
	if err != nil {
		return nil, fmt.Errorf("parse balance: %w", err)
	}
	return &domain.Account{ID: id, Balance: bal, UpdatedAt: t}, nil
}

func (r *BalanceRepository) ProcessTransaction(ctx context.Context, op *domain.Operation) (*domain.Account, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback()

	// creating account
	if _, err := tx.ExecContext(ctx, sqlCreateAccount, op.AccountID); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	res, err := tx.ExecContext(ctx, sqlInsertOperation,
		op.TxID, op.AccountID, string(op.Source), string(op.State), op.Amount.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("insert op: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		// if operation exist - just return current balance
		acc, err := selectAccount(ctx, tx, op.AccountID)
		if err != nil {
			return nil, fmt.Errorf("select balance (dup): %w", err)
		}
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit (dup): %w", err)
		}
		return acc, domain.ErrDuplicateTx
	}

	// compute delta for the balance update
	delta := op.Amount
	if op.State == domain.StateWithdraw {
		delta = delta.Neg()
	}

	// appply delta with non-negative constraint
	var s string
	var t time.Time
	if err := tx.QueryRowContext(ctx, sqlUpdateBalance, delta.String(), op.AccountID).Scan(&s, &t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNegativeBalance
		}
		return nil, fmt.Errorf("update balance: %w", err)
	}
	bal, err := decimal.NewFromString(s)
	if err != nil {
		return nil, fmt.Errorf("parse updated balance: %w", err)
	}

	// operation sucessful
	if _, err := tx.ExecContext(ctx, `UPDATE operations SET applied = true WHERE tx_id = $1`, op.TxID); err != nil {
		return nil, fmt.Errorf("mark applied: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &domain.Account{ID: op.AccountID, Balance: bal, UpdatedAt: t}, nil
}
