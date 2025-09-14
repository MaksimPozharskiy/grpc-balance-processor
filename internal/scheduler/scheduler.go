package scheduler

import (
	"context"
	"time"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/db"
	"go.uber.org/zap"
)

const AdvisoryLockKey = int64(0xBABACAFE)

type Scheduler struct {
	db     *db.DB
	period time.Duration
	log    *zap.Logger
}

func NewCancelScheduler(database *db.DB, period time.Duration, log *zap.Logger) *Scheduler {
	return &Scheduler{
		db:     database,
		period: period,
		log:    log.Named("cancel-scheduler"),
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.log.Info("starting cancel scheduler", zap.Duration("period", s.period))

	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	s.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			s.log.Info("cancel scheduler stopped")
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	s.log.Debug("starting cancellation cycle")

	if !s.tryAdvisoryLock(ctx) {
		s.log.Debug("cancel scheduler: lock busy, skipping cycle")
		return
	}
	defer s.advisoryUnlock(ctx)

	candidates, err := s.selectOddLatestOps(ctx)
	if err != nil {
		s.log.Error("failed to select candidate operations", zap.Error(err))
		return
	}

	s.log.Info("starting cancellation cycle", zap.Int("candidates", len(candidates)))

	var cancelled, skipped, failed int

	for _, candidate := range candidates {
		switch result := s.cancelOne(ctx, candidate.ID); result {
		case CancelResultSuccess:
			cancelled++
			s.log.Info("operation cancelled successfully",
				zap.Int64("op_id", candidate.ID),
				zap.String("tx_id", candidate.TxID))
		case CancelResultSkipped:
			skipped++
			s.log.Info("operation skipped due to insufficient funds",
				zap.Int64("op_id", candidate.ID),
				zap.String("tx_id", candidate.TxID))
		case CancelResultFailed:
			failed++
			s.log.Warn("operation cancellation failed",
				zap.Int64("op_id", candidate.ID),
				zap.String("tx_id", candidate.TxID))
		}
	}

	s.log.Info("cancellation cycle completed",
		zap.Int("total_candidates", len(candidates)),
		zap.Int("cancelled", cancelled),
		zap.Int("skipped", skipped),
		zap.Int("failed", failed))
}
