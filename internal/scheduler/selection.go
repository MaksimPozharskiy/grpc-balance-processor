package scheduler

import (
	"context"

	"go.uber.org/zap"
)

func (s *Scheduler) selectOddLatestOps(ctx context.Context) ([]CandidateOperation, error) {
	query := `
		WITH ranked AS (
			SELECT id, tx_id
			FROM (
				SELECT o.*,
					ROW_NUMBER() OVER (ORDER BY o.created_at DESC, o.id DESC) AS rn
				FROM operations o
				WHERE o.applied = TRUE
					AND o.canceled_at IS NULL
			) t
			WHERE (rn % 2) = 1
			ORDER BY rn
			LIMIT 10
		)
		SELECT id, tx_id
		FROM ranked;`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		s.log.Error("failed to query odd operations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var candidates []CandidateOperation
	for rows.Next() {
		var candidate CandidateOperation
		if err := rows.Scan(&candidate.ID, &candidate.TxID); err != nil {
			s.log.Error("failed to scan candidate operation", zap.Error(err))
			return nil, err
		}
		candidates = append(candidates, candidate)
	}

	if err := rows.Err(); err != nil {
		s.log.Error("error iterating candidate operations", zap.Error(err))
		return nil, err
	}

	s.log.Debug("selected candidate operations", zap.Int("count", len(candidates)))
	return candidates, nil
}
