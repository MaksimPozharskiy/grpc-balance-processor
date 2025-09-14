package config

import (
	"os"
	"testing"
)

func TestConfig_Defaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_DSN", "test-dsn")

	cfg := Load()

	if cfg.GRPCPort != "8080" {
		t.Errorf("GRPCPort = %v, want 8080", cfg.GRPCPort)
	}

	if cfg.CancelPeriodMin != 5 {
		t.Errorf("CancelPeriodMin = %v, want 5", cfg.CancelPeriodMin)
	}

	if cfg.CancelSchedulerEnabled != true {
		t.Errorf("CancelSchedulerEnabled = %v, want true", cfg.CancelSchedulerEnabled)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %v, want info", cfg.LogLevel)
	}

	if cfg.DatabaseDSN != "test-dsn" {
		t.Errorf("DatabaseDSN = %v, want test-dsn", cfg.DatabaseDSN)
	}
}

func TestConfig_CustomValues(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_DSN", "custom-dsn")
	os.Setenv("GRPC_PORT", "9090")
	os.Setenv("CANCEL_PERIOD_MIN", "10")
	os.Setenv("CANCEL_SCHEDULER_ENABLED", "false")
	os.Setenv("LOG_LEVEL", "debug")

	cfg := Load()

	if cfg.GRPCPort != "9090" {
		t.Errorf("GRPCPort = %v, want 9090", cfg.GRPCPort)
	}

	if cfg.CancelPeriodMin != 10 {
		t.Errorf("CancelPeriodMin = %v, want 10", cfg.CancelPeriodMin)
	}

	if cfg.CancelSchedulerEnabled != false {
		t.Errorf("CancelSchedulerEnabled = %v, want false", cfg.CancelSchedulerEnabled)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %v, want debug", cfg.LogLevel)
	}
}
