package transport

import (
	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	pb "github.com/MaksimPozharskiy/grpc-balance-processor/proto/balance"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
