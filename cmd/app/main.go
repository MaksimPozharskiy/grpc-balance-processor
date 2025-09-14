package main

import (
	"context"
	"time"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/config"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/db"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/logger"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/repository"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/scheduler"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/transport"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "balance-service", "dev", "local")
	logger.SetGlobal(log)
	defer logger.Sync()

	log.Info("starting application",
		zap.String("grpc_port", cfg.GRPCPort),
		zap.Bool("cancel_scheduler_enabled", cfg.CancelSchedulerEnabled),
		zap.Int("cancel_period_min", cfg.CancelPeriodMin),
	)

	database, err := db.NewConnection(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.HealthCheck(ctx); err != nil {
		log.Fatal("database health check failed", zap.Error(err))
	}

	repo := repository.NewBalanceRepository(database)
	balanceService := usecase.NewBalanceUsecase(repo)

	if cfg.CancelSchedulerEnabled {
		cancelScheduler := scheduler.NewCancelScheduler(
			database,
			time.Duration(cfg.CancelPeriodMin)*time.Minute,
			log,
		)
		go cancelScheduler.Run(ctx)
		log.Info("cancel scheduler started")
	} else {
		log.Info("cancel scheduler disabled")
	}

	server := transport.NewGRPCServer(balanceService)

	if err := transport.Serve(server, cfg.GRPCPort); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
