package main

import (
	"context"
	"net"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/config"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/db"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/logger"
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/repository"
	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedBalanceServiceServer
}

func (s *server) Process(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	// TODO заглушка
	return &pb.ProcessResponse{
		TxId:        req.TxId,
		Status:      pb.Status_STATUS_OK,
		Balance:     "100.00",
		ProcessedAt: timestamppb.Now(),
	}, nil
}

func (s *server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	// TODO заглушка
	return &pb.GetBalanceResponse{
		Balance:   "100.00",
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "balance-service", "dev", "local")
	logger.SetGlobal(log)
	defer logger.Sync()

	log.Info("starting application",
		zap.String("grpc_port", cfg.GRPCPort),
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
	_ = repo // TODO доделать как слои появятся другие

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

	pb.RegisterBalanceServiceServer(s, &server{})

	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthService)
	healthService.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	log.Info("gRPC server started", zap.String("port", cfg.GRPCPort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
