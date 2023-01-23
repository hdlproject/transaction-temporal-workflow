package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
	"transaction-temporal-workflow/usecase"
)

type (
	UseCase struct {
		transactionCommand repository.TransactionCommand
		transactionQuery   repository.TransactionQuery
		userRepository     repository.User
		productRepository  repository.Product

		rabbitMQ *amqp.Channel
	}

	CreatedTransaction struct {
		UseCase
	}

	PendingTransaction struct {
		UseCase
	}
)

func NewUseCase(transactionRepository repository.Transaction, userRepository repository.User, productRepository repository.Product, rabbitMQ *amqp.Channel) UseCase {
	return UseCase{
		transactionCommand: transactionRepository.Command,
		transactionQuery:   transactionRepository.Query,
		userRepository:     userRepository,
		productRepository:  productRepository,
		rabbitMQ:           rabbitMQ,
	}
}

func (i UseCase) CreateTransaction(ctx context.Context, transaction model.Transaction) error {
	_, err := i.transactionQuery.GetLastTransactionByTransactionId(transaction.TransactionId)
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
		IsPublished:   false,
	}

	err = i.transactionCommand.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (i UseCase) ProcessTransaction(ctx context.Context, transactionId string, expectedStatus model.TransactionStatus) error {
	transaction, err := i.transactionQuery.GetLastTransactionByTransactionId(transactionId)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	switch transaction.Status {
	case model.TransactionStatusCreated:
		return CreatedTransaction{i}.ProcessTransaction(transaction, expectedStatus)
	case model.TransactionStatusPending:
		return PendingTransaction{i}.ProcessTransaction(transaction, expectedStatus)
	default:
		return fmt.Errorf("transaction status %s is already final", transaction.Status)
	}
}

func (i CreatedTransaction) ProcessTransaction(transaction model.Transaction, expectedStatus model.TransactionStatus) error {
	var nextStatus model.TransactionStatus
	switch expectedStatus {
	case model.TransactionStatusPending, model.TransactionStatusFailed:
		nextStatus = expectedStatus
	default:
		return fmt.Errorf("transction status transition from %s to %s is not allowed", transaction.Status, expectedStatus)
	}

	transaction.Id = 0
	transaction.Status = nextStatus
	transaction.CreatedAt = time.Now()

	err := i.transactionCommand.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("create transaction %w", err)
	}

	return nil
}

func (i PendingTransaction) ProcessTransaction(transaction model.Transaction, expectedStatus model.TransactionStatus) error {
	var nextStatus model.TransactionStatus
	switch expectedStatus {
	case model.TransactionStatusSuccess, model.TransactionStatusFailed:
		nextStatus = expectedStatus
	default:
		return fmt.Errorf("transction status transition from %s to %s is not allowed", transaction.Status, expectedStatus)
	}

	transaction.Id = 0
	transaction.Status = nextStatus
	transaction.CreatedAt = time.Now()

	err := i.transactionCommand.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("create transaction %w", err)
	}

	return nil
}

func (i UseCase) PublishTransaction(ctx context.Context) error {
	transactions, err := i.transactionQuery.GetUnpublishedTransactions()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("get unpublished transactions: %w", err)
		}

		return nil
	}

	for _, transaction := range transactions {
		transactionJson, err := json.Marshal(transaction)
		if err != nil {
			return fmt.Errorf("json marshall: %w", err)
		}

		err = i.rabbitMQ.PublishWithContext(ctx,
			usecase.TransactionExchangeName,
			usecase.TransactionCreatedRoutingKey,
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

		err = i.transactionCommand.PublishTransaction(transaction.Id)
		if err != nil {
			return fmt.Errorf("publish transaction: %w", err)
		}
	}

	return nil
}
