package transport

import (
	"context"
	"net"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
		return nil, err
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
		return nil, err
	}

	return &pb.GetBalanceResponse{
		Balance:   resp.Balance.String(),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

// TODO надо вынести отдельно
func mapProtoSource(s pb.Source) (domain.Source, error) {
	switch s {
	case pb.Source_SOURCE_GAME:
		return domain.SourceGame, nil
	case pb.Source_SOURCE_PAYMENT:
		return domain.SourcePayment, nil
	case pb.Source_SOURCE_SERVICE:
		return domain.SourceService, nil
	default:
		return "", status.Error(codes.InvalidArgument, "invalid source value")
	}
}

func mapProtoState(s pb.State) (domain.State, error) {
	switch s {
	case pb.State_STATE_DEPOSIT:
		return domain.StateDeposit, nil
	case pb.State_STATE_WITHDRAW:
		return domain.StateWithdraw, nil
	default:
		return "", status.Error(codes.InvalidArgument, "invalid state value")
	}
}

func mapDomainStatus(s domain.ProcessStatus) pb.Status {
	switch s {
	case domain.StatusOK:
		return pb.Status_STATUS_OK
	case domain.StatusAlreadyProcessed:
		return pb.Status_STATUS_ALREADY_PROCESSED
	case domain.StatusRejectedNegative:
		return pb.Status_STATUS_REJECTED_NEGATIVE
	default:
		return pb.Status_STATUS_OK
	}
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
