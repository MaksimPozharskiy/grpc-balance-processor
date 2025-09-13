package usecase

import (
	"context"
	"errors"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"go.uber.org/zap"
)

type BalanceUsecase struct {
	repo domain.BalanceRepository
}

func NewBalanceUsecase(repo domain.BalanceRepository) *BalanceUsecase {
	return &BalanceUsecase{
		repo: repo,
	}
}

func (u *BalanceUsecase) Process(ctx context.Context, req *domain.ProcessRequest) (*domain.ProcessResponse, error) {
	zap.L().Info("processing transaction",
		zap.String("tx_id", req.TxID),
		zap.String("account_id", req.AccountID.String()),
		zap.String("amount", req.Amount.String()),
	)

	op := &domain.Operation{
		TxID:      req.TxID,
		AccountID: req.AccountID,
		Source:    req.Source,
		State:     req.State,
		Amount:    req.Amount,
	}

	account, err := u.repo.ProcessTransaction(ctx, op)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateTx) {
			return &domain.ProcessResponse{
				TxID:      req.TxID,
				Status:    domain.StatusAlreadyProcessed,
				Balance:   account.Balance,
				Timestamp: account.UpdatedAt,
			}, nil
		}
		if errors.Is(err, domain.ErrNegativeBalance) {
			currentAccount, getErr := u.repo.GetAccount(ctx, req.AccountID)
			if getErr != nil {
				return nil, getErr
			}
			return &domain.ProcessResponse{
				TxID:      req.TxID,
				Status:    domain.StatusRejectedNegative,
				Balance:   currentAccount.Balance,
				Timestamp: currentAccount.UpdatedAt,
			}, nil
		}
		zap.L().Error("ProcessTransaction failed", zap.Error(err))
		return nil, err
	}

	return &domain.ProcessResponse{
		TxID:      req.TxID,
		Status:    domain.StatusOK,
		Balance:   account.Balance,
		Timestamp: account.UpdatedAt,
	}, nil
}

func (u *BalanceUsecase) GetBalance(ctx context.Context, req *domain.GetBalanceRequest) (*domain.GetBalanceResponse, error) {
	zap.L().Info("getting balance",
		zap.String("account_id", req.AccountID.String()),
	)

	account, err := u.repo.GetAccount(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &domain.GetBalanceResponse{
		Balance:   account.Balance,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
