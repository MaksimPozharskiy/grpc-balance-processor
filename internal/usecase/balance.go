package usecase

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
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

	// заглушка
	return &domain.ProcessResponse{
		TxID:      req.TxID,
		Status:    domain.StatusOK,
		Balance:   decimal.NewFromFloat(100.00),
		Timestamp: time.Now(),
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
