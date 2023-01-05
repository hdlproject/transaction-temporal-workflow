package activity

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
)

type (
	Transaction struct {
		transactionRepository repository.Transaction
	}

	CreatedTransaction struct {
		Transaction
	}

	PendingTransaction struct {
		Transaction
	}
)

func NewTransaction(transactionRepository repository.Transaction) Transaction {
	return Transaction{
		transactionRepository: transactionRepository,
	}
}

func (i Transaction) ProcessTransaction(transactionId, idempotencyKey string) error {
	isAllowed, err := i.transactionRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	transaction, err := i.transactionRepository.GetTransaction(transactionId)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err != nil && err == redis.Nil {
		err = i.transactionRepository.CreateTransaction(transactionId)
		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		return nil
	}

	switch transaction.Status {
	case model.TransactionStatusCreated:
		return CreatedTransaction{i}.ProcessTransaction(transactionId)
	case model.TransactionStatusPending:
		return PendingTransaction{i}.ProcessTransaction(transactionId)
	}

	return nil
}

func (i CreatedTransaction) ProcessTransaction(transactionId string) error {
	return i.transactionRepository.UpdateTransactionStatus(transactionId, model.TransactionStatusPending)
}

func (i PendingTransaction) ProcessTransaction(transactionId string) error {
	return i.transactionRepository.UpdateTransactionStatus(transactionId, model.TransactionStatusSuccess)
}
