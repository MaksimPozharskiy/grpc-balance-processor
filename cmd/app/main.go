package main

import (
	"context"
	"log"
	"net"

	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
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
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterBalanceServiceServer(s, &server{})

	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthService)
	healthService.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	log.Println("gRPC server listening on :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
