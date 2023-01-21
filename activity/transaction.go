package activity

import (
	"context"
	"fmt"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/usecase/transaction"
)

type (
	Transaction struct {
		transactionUseCase transaction.UseCase
	}
)

func NewTransaction(transactionUseCase transaction.UseCase) Transaction {
	return Transaction{
		transactionUseCase: transactionUseCase,
	}
}

func (i Transaction) CreateTransaction(ctx context.Context, transaction model.Transaction) error {
	err := i.transactionUseCase.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (i Transaction) ProcessTransaction(ctx context.Context, transactionId string, expectedStatus model.TransactionStatus) error {
	err := i.transactionUseCase.ProcessTransaction(ctx, transactionId, expectedStatus)
	if err != nil {
		return fmt.Errorf("process transaction: %w", err)
	}

	return nil
}

func (i Transaction) PublishTransaction(ctx context.Context) error {
	err := i.transactionUseCase.PublishTransaction(ctx)
	if err != nil {
		return fmt.Errorf("publish transaction: %w", err)
	}

	return nil
}
