package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
)

type (
	UseCase struct {
		transactionCommand    repository.TransactionCommand
		transactionQuery      repository.TransactionQuery
		idempotencyRepository repository.Idempotency

		rabbitMQ *amqp.Channel
	}

	CreatedTransaction struct {
		UseCase
	}

	PendingTransaction struct {
		UseCase
	}
)

func NewUseCase(transactionRepository repository.Transaction, idempotencyRepository repository.Idempotency, rabbitMQ *amqp.Channel) UseCase {
	return UseCase{
		transactionCommand:    transactionRepository.Command,
		transactionQuery:      transactionRepository.Query,
		idempotencyRepository: idempotencyRepository,
		rabbitMQ:              rabbitMQ,
	}
}

func (i UseCase) CreateTransaction(ctx context.Context, transaction model.Transaction, idempotencyKey string) error {
	isAllowed, err := i.idempotencyRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

	_, err = i.transactionQuery.GetLastTransactionByTransactionId(transaction.TransactionId)
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

	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("json marshall: %w", err)
	}

	err = i.rabbitMQ.PublishWithContext(ctx,
		cmd.TransactionExchangeName,
		cmd.TransactionCreatedRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        transactionJson,
		},
	)
	if err != nil {
		return fmt.Errorf("rabbitmq publish with context: %w", err)
	}

	return nil
}

func (i UseCase) ProcessTransaction(ctx context.Context, transactionId, idempotencyKey string) error {
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
