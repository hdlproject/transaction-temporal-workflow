package activity

import (
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

func (i Transaction) CreateTransaction(transaction model.Transaction, idempotencyKey string) error {
	err := i.transactionUseCase.CreateTransaction(transaction, idempotencyKey)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (i Transaction) ProcessTransaction(transactionId, idempotencyKey string) error {
	err := i.transactionUseCase.ProcessTransaction(transactionId, idempotencyKey)
	if err != nil {
		return fmt.Errorf("process transaction: %w", err)
	}

	return nil
}
