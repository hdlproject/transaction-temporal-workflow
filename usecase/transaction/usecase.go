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

	product, err := i.productRepository.GetProductByCode(transaction.ProductCode)
	if err != nil {
		return fmt.Errorf("get product by code: %w", err)
	}

	user, err := i.userRepository.GetUserById(transaction.UserId)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	transaction = model.Transaction{
		Id:            transaction.Id,
		TransactionId: transaction.TransactionId,
		Status:        model.TransactionStatusCreated,
		Amount:        transaction.Amount,
		ProductCode:   transaction.ProductCode,
		Product:       product,
		UserId:        transaction.UserId,
		User:          user,
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

	return nil
}

func (i UseCase) ProcessTransaction(ctx context.Context, transactionId string) error {
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
