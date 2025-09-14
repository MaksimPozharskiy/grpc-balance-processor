package domain

import (
	"testing"
)

func TestSource_Constants(t *testing.T) {
	tests := []struct {
		source   Source
		expected string
	}{
		{SourceGame, "game"},
		{SourcePayment, "payment"},
		{SourceService, "service"},
	}

	for _, tt := range tests {
		if string(tt.source) != tt.expected {
			t.Errorf("Source %v = %v, want %v", tt.source, string(tt.source), tt.expected)
		}
	}
}

func TestState_Constants(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateDeposit, "deposit"},
		{StateWithdraw, "withdraw"},
	}

	for _, tt := range tests {
		if string(tt.state) != tt.expected {
			t.Errorf("State %v = %v, want %v", tt.state, string(tt.state), tt.expected)
		}
	}
}

func TestProcessStatus_Constants(t *testing.T) {
	if StatusOK != 0 {
		t.Errorf("StatusOK = %v, want 0", StatusOK)
	}
	if StatusAlreadyProcessed != 1 {
		t.Errorf("StatusAlreadyProcessed = %v, want 1", StatusAlreadyProcessed)
	}
	if StatusRejectedNegative != 2 {
		t.Errorf("StatusRejectedNegative = %v, want 2", StatusRejectedNegative)
	}
}