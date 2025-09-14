package scheduler

import (
	"context"
)

type Canceler interface {
	Run(ctx context.Context)
}

type CandidateOperation struct {
	ID   int64
	TxID string
}
