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
		transactionRepository repository.Transaction
	}

	CreatedTransaction struct {
		UseCase
	}

	PendingTransaction struct {
		UseCase
	}
)

func NewUseCase(transactionRepository repository.Transaction) UseCase {
	return UseCase{
		transactionRepository: transactionRepository,
	}
}

func (i UseCase) CreateTransaction(transaction model.Transaction, idempotencyKey string) error {
	isAllowed, err := i.transactionRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	_, err = i.transactionRepository.GetTransactionByTransactionId(transaction.TransactionId)
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

	err = i.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (i UseCase) ProcessTransaction(transactionId, idempotencyKey string) error {
	isAllowed, err := i.transactionRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	transaction, err := i.transactionRepository.GetLastTransactionByTransactionId(transactionId)
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

	return i.transactionRepository.CreateTransaction(transaction)
}

func (i PendingTransaction) ProcessTransaction(transaction model.Transaction) error {
	transaction.Id = 0
	transaction.Status = model.TransactionStatusSuccess
	transaction.CreatedAt = time.Now()

	return i.transactionRepository.CreateTransaction(transaction)
}
