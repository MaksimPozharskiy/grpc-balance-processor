package scheduler

import (
	"testing"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"github.com/shopspring/decimal"
)

func TestCalculateCompensatingDelta(t *testing.T) {
	s := &Scheduler{}

	tests := []struct {
		state    domain.State
		amount   decimal.Decimal
		expected decimal.Decimal
	}{
		{domain.StateDeposit, decimal.NewFromFloat(10.50), decimal.NewFromFloat(-10.50)},
		{domain.StateWithdraw, decimal.NewFromFloat(5.25), decimal.NewFromFloat(5.25)},
	}

	for _, tt := range tests {
		result := s.calculateCompensatingDelta(tt.state, tt.amount)
		if !result.Equal(tt.expected) {
			t.Errorf("calculateCompensatingDelta(%v, %v) = %v, want %v",
				tt.state, tt.amount, result, tt.expected)
		}
	}
}

func TestGetCompensatingState(t *testing.T) {
	s := &Scheduler{}

	tests := []struct {
		original domain.State
		expected domain.State
	}{
		{domain.StateDeposit, domain.StateWithdraw},
		{domain.StateWithdraw, domain.StateDeposit},
	}

	for _, tt := range tests {
		result := s.getCompensatingState(tt.original)
		if result != tt.expected {
			t.Errorf("getCompensatingState(%v) = %v, want %v",
				tt.original, result, tt.expected)
		}
	}
}
