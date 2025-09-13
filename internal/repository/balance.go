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

type BalanceRepository struct {
	db *db.DB
}

func NewBalanceRepository(database *db.DB) *BalanceRepository {
	return &BalanceRepository{db: database}
}

func (r *BalanceRepository) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	query := `SELECT id, balance, updated_at FROM accounts WHERE id = $1`

	var acc domain.Account
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&acc.ID, &acc.Balance, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get account: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get account query: %w", err)
	}

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

func (r *BalanceRepository) ProcessTransaction(ctx context.Context, op *domain.Operation) (*domain.Account, error) {
	// TODO тут транзакции попомзже напишу
	return &domain.Account{
		ID:        op.AccountID,
		Balance:   decimal.NewFromFloat(100.00),
		UpdatedAt: time.Now(),
	}, nil
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
