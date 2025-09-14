package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/MaksimPozharskiy/grpc-balance-processor/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	account   *domain.Account
	operation *domain.Operation
	err       error
}

func (m *mockRepository) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	return m.account, m.err
}

func (m *mockRepository) CreateAccount(ctx context.Context, accountID uuid.UUID) error {
	return m.err
}

func (m *mockRepository) ProcessTransaction(ctx context.Context, op *domain.Operation) (*domain.Account, error) {
	return m.account, m.err
}

func (m *mockRepository) GetOperationByTxID(ctx context.Context, txID string) (*domain.Operation, error) {
	return m.operation, m.err
}

func TestBalanceUsecase_Process_Deposit(t *testing.T) {
	accountID := uuid.New()
	mockRepo := &mockRepository{
		account: &domain.Account{
			ID:        accountID,
			Balance:   decimal.NewFromFloat(110.50),
			UpdatedAt: time.Now(),
		},
	}

	usecase := NewBalanceUsecase(mockRepo)

	req := &domain.ProcessRequest{
		AccountID: accountID,
		Source:    domain.SourceGame,
		State:     domain.StateDeposit,
		Amount:    decimal.NewFromFloat(10.50),
		TxID:      "test-tx-001",
	}

	resp, err := usecase.Process(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusOK, resp.Status)
	assert.True(t, resp.Balance.Equal(decimal.NewFromFloat(110.50)))
}

func TestBalanceUsecase_GetBalance(t *testing.T) {
	accountID := uuid.New()
	mockRepo := &mockRepository{
		account: &domain.Account{
			ID:        accountID,
			Balance:   decimal.NewFromFloat(25.75),
			UpdatedAt: time.Now(),
		},
	}

	usecase := NewBalanceUsecase(mockRepo)

	req := &domain.GetBalanceRequest{
		AccountID: accountID,
	}

	resp, err := usecase.GetBalance(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, resp.Balance.Equal(decimal.NewFromFloat(25.75)))
}
