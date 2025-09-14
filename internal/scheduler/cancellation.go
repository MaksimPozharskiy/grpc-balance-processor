package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type CancelResult int

const (
	CancelResultSuccess CancelResult = iota
	CancelResultSkipped
	CancelResultFailed
)

func (s *Scheduler) cancelOne(ctx context.Context, operationID int64) CancelResult {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Error("failed to begin transaction for cancellation", zap.Int64("op_id", operationID), zap.Error(err))
		return CancelResultFailed
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.log.Error("failed to rollback transaction", zap.Error(err))
		}
	}()

	// load operation
	operation, err := s.loadOperation(ctx, tx, operationID)
	if err != nil {
		s.log.Error("failed to load operation for cancellation", zap.Int64("op_id", operationID), zap.Error(err))
		return CancelResultFailed
	}

	// idempotency: skip if not applied or already canceled
	if !operation.Applied || operation.CanceledAt != nil {
		s.log.Debug("operation already cancelled or not applied", zap.Int64("op_id", operationID))
		if commitErr := tx.Commit(); commitErr != nil {
			s.log.Error("failed to commit transaction", zap.Error(commitErr))
		}
		return CancelResultSkipped
	}

	// compute delta
	compensatingDelta := s.calculateCompensatingDelta(operation.State, operation.Amount)

	// apply delta with non-negative guard
	balanceUpdated, err := s.updateAccountBalance(ctx, tx, operation.AccountID, compensatingDelta)
	if err != nil {
		s.log.Error("failed to update account balance", zap.Int64("op_id", operationID), zap.Error(err))
		return CancelResultFailed
	}

	// mark skipped and commit
	if !balanceUpdated {
		if err := s.markOperationAsSkipped(ctx, tx, operationID); err != nil {
			s.log.Error("failed to mark operation as skipped", zap.Int64("op_id", operationID), zap.Error(err))
			return CancelResultFailed
		}

		if err := tx.Commit(); err != nil {
			s.log.Error("failed to commit skip transaction", zap.Error(err))
			return CancelResultFailed
		}
		return CancelResultSkipped
	}

	compensatingTxID := fmt.Sprintf("cancel::%s", operation.TxID)

	// write compensating operation
	if err := s.createCompensatingOperation(ctx, tx, operation, compensatingTxID, compensatingDelta); err != nil {
		s.log.Error("failed to create compensating operation", zap.Int64("op_id", operationID), zap.Error(err))
		return CancelResultFailed
	}

	// / mark original as canceled
	if err := s.markOperationAsCancelled(ctx, tx, operationID); err != nil {
		s.log.Error("failed to mark operation as cancelled", zap.Int64("op_id", operationID), zap.Error(err))
		return CancelResultFailed
	}

	// commit
	if err := tx.Commit(); err != nil {
		s.log.Error("failed to commit cancellation transaction", zap.Error(err))
		return CancelResultFailed
	}

	return CancelResultSuccess
}

func (s *Scheduler) loadOperation(ctx context.Context, tx *sql.Tx, operationID int64) (*domain.Operation, error) {
	query := `
		SELECT id, tx_id, account_id, source, state, amount, created_at, applied, canceled_at, cancel_note
		FROM operations
		WHERE id = $1`

	var op domain.Operation
	var canceledAt *time.Time
	var cancelNote *string

	err := tx.QueryRowContext(ctx, query, operationID).Scan(
		&op.ID, &op.TxID, &op.AccountID, &op.Source, &op.State, &op.Amount,
		&op.CreatedAt, &op.Applied, &canceledAt, &cancelNote,
	)
	if err != nil {
		return nil, err
	}

	op.CanceledAt = canceledAt
	op.CancelNote = cancelNote
	return &op, nil
}

func (s *Scheduler) calculateCompensatingDelta(state domain.State, amount decimal.Decimal) decimal.Decimal {
	switch state {
	case domain.StateDeposit:
		return amount.Neg()
	case domain.StateWithdraw:
		return amount
	default:
		return decimal.Zero
	}
}

func (s *Scheduler) updateAccountBalance(ctx context.Context, tx *sql.Tx, accountID uuid.UUID, delta decimal.Decimal) (bool, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = now()
		WHERE id = $2 AND balance + $1 >= 0
		RETURNING balance, updated_at`

	var balance decimal.Decimal
	var updatedAt time.Time
	err := tx.QueryRowContext(ctx, query, delta, accountID).Scan(&balance, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Scheduler) createCompensatingOperation(ctx context.Context, tx *sql.Tx, originalOp *domain.Operation, compensatingTxID string, amount decimal.Decimal) error {
	compensatingState := s.getCompensatingState(originalOp.State)

	query := `
		INSERT INTO operations (tx_id, account_id, source, state, amount, applied, cancel_note)
		VALUES ($1, $2, $3, $4, $5, TRUE, 'auto-cancel')
		ON CONFLICT (tx_id) DO NOTHING`

	_, err := tx.ExecContext(ctx, query,
		compensatingTxID,
		originalOp.AccountID,
		domain.SourceService,
		compensatingState,
		amount.Abs(),
	)
	return err
}

func (s *Scheduler) getCompensatingState(originalState domain.State) domain.State {
	switch originalState {
	case domain.StateDeposit:
		return domain.StateWithdraw
	case domain.StateWithdraw:
		return domain.StateDeposit
	default:
		return domain.StateDeposit
	}
}

func (s *Scheduler) markOperationAsCancelled(ctx context.Context, tx *sql.Tx, operationID int64) error {
	query := `
		UPDATE operations
		SET canceled_at = now(), cancel_note = 'scheduler'
		WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, operationID)
	return err
}

func (s *Scheduler) markOperationAsSkipped(ctx context.Context, tx *sql.Tx, operationID int64) error {
	query := `
		UPDATE operations
		SET cancel_note = 'skip: insufficient funds'
		WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, operationID)
	return err
}
