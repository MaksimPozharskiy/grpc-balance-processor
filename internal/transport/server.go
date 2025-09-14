package transport

import (
	"context"
	"net"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedBalanceServiceServer
	service domain.BalanceService
}

func NewServer(service domain.BalanceService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Process(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	if err := validateProcessRequest(req); err != nil {
		return nil, err
	}

	accountID, err := validateAndParseAccountID(req.AccountId)
	if err != nil {
		return nil, err
	}

	amount, err := validateAndParseAmount(req.Amount)
	if err != nil {
		return nil, err
	}

	if err := validateTxID(req.TxId); err != nil {
		return nil, err
	}

	source, err := mapProtoSource(req.Source)
	if err != nil {
		return nil, err
	}

	state, err := mapProtoState(req.State)
	if err != nil {
		return nil, err
	}

	domainReq := &domain.ProcessRequest{
		AccountID: accountID,
		Source:    source,
		State:     state,
		Amount:    amount,
		TxID:      req.TxId,
	}

	resp, err := s.service.Process(ctx, domainReq)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.ProcessResponse{
		TxId:        resp.TxID,
		Status:      mapDomainStatus(resp.Status),
		Balance:     resp.Balance.String(),
		ProcessedAt: timestamppb.New(resp.Timestamp),
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	accountID, err := validateAndParseAccountID(req.AccountId)
	if err != nil {
		return nil, err
	}

	domainReq := &domain.GetBalanceRequest{
		AccountID: accountID,
	}

	resp, err := s.service.GetBalance(ctx, domainReq)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.GetBalanceResponse{
		Balance:   resp.Balance.String(),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

func NewGRPCServer(service domain.BalanceService) *grpc.Server {
	s := grpc.NewServer()

	pb.RegisterBalanceServiceServer(s, NewServer(service))

	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthService)
	healthService.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	return s
}

func Serve(s *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	zap.L().Info("gRPC server started", zap.String("port", port))
	return s.Serve(lis)
}
