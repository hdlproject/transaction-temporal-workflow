package transaction

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
)

type (
	UseCase struct {
		transactionCommand    repository.TransactionCommand
		transactionQuery      repository.TransactionQuery
		idempotencyRepository repository.Idempotency
	}

	CreatedTransaction struct {
		UseCase
	}

	PendingTransaction struct {
		UseCase
	}
)

func NewUseCase(transactionRepository repository.Transaction, idempotencyRepository repository.Idempotency) UseCase {
	return UseCase{
		transactionCommand:    transactionRepository.Command,
		transactionQuery:      transactionRepository.Query,
		idempotencyRepository: idempotencyRepository,
	}
}

func (i UseCase) CreateTransaction(transaction model.Transaction, idempotencyKey string) error {
	isAllowed, err := i.idempotencyRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	_, err = i.transactionQuery.GetTransactionByTransactionId(transaction.TransactionId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err == nil {
		return fmt.Errorf("record already exists")
	}

	transaction = model.Transaction{
		Id:            transaction.Id,
		TransactionId: transaction.TransactionId,
		Status:        model.TransactionStatusCreated,
		Amount:        transaction.Amount,
		ProductCode:   transaction.ProductCode,
		UserId:        transaction.UserId,
		CreatedAt:     time.Now(),
	}

	err = i.transactionCommand.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (i UseCase) ProcessTransaction(transactionId, idempotencyKey string) error {
	isAllowed, err := i.idempotencyRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	transaction, err := i.transactionQuery.GetLastTransactionByTransactionId(transactionId)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	switch transaction.Status {
	case model.TransactionStatusCreated:
		return CreatedTransaction{i}.ProcessTransaction(transaction)
	case model.TransactionStatusPending:
		return PendingTransaction{i}.ProcessTransaction(transaction)
	}

	return nil
}

func (i CreatedTransaction) ProcessTransaction(transaction model.Transaction) error {
	transaction.Id = 0
	transaction.Status = model.TransactionStatusPending
	transaction.CreatedAt = time.Now()

	return i.transactionCommand.CreateTransaction(transaction)
}

func (i PendingTransaction) ProcessTransaction(transaction model.Transaction) error {
	transaction.Id = 0
	transaction.Status = model.TransactionStatusSuccess
	transaction.CreatedAt = time.Now()

	return i.transactionCommand.CreateTransaction(transaction)
}
