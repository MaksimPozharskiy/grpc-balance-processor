package scheduler

import (
	"context"

	"go.uber.org/zap"
)

func (s *Scheduler) tryAdvisoryLock(ctx context.Context) bool {
	var acquired bool
	err := s.db.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", AdvisoryLockKey).Scan(&acquired)
	if err != nil {
		s.log.Error("failed to acquire advisory lock", zap.Error(err))
		return false
	}

	s.log.Debug("advisory lock attempt", zap.Bool("acquired", acquired))
	return acquired
}

func (s *Scheduler) advisoryUnlock(ctx context.Context) {
	var released bool
	err := s.db.QueryRowContext(ctx, "SELECT pg_advisory_unlock($1)", AdvisoryLockKey).Scan(&released)
	if err != nil {
		s.log.Error("failed to release advisory lock", zap.Error(err))
		return
	}

	s.log.Debug("advisory lock released", zap.Bool("released", released))
}
