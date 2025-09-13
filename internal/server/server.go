package server

import (
	"context"

	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedBalanceServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Process(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	return &pb.ProcessResponse{
		TxId:        req.TxId,
		Status:      pb.Status_STATUS_OK,
		Balance:     "100.00",
		ProcessedAt: timestamppb.Now(),
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	return &pb.GetBalanceResponse{
		Balance:   "100.00",
		UpdatedAt: timestamppb.Now(),
	}, nil
}
