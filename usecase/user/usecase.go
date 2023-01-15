package user

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
	"transaction-temporal-workflow/usecase"
	"transaction-temporal-workflow/usecase/idempotency"
)

type (
	UseCase struct {
		idempotencyUseCase idempotency.UseCase

		userRepository repository.User

		rabbitMQ *amqp.Channel
	}
)

func NewUseCase(idempotencyUseCase idempotency.UseCase, userRepository repository.User, rabbitMQ *amqp.Channel) UseCase {
	return UseCase{
		idempotencyUseCase: idempotencyUseCase,
		userRepository:     userRepository,
		rabbitMQ:           rabbitMQ,
	}
}

func (i UseCase) ReserveUserBalance(ctx context.Context, transaction model.Transaction) error {
	var transactionStatus model.TransactionStatus
	err := i.reserveUserBalance(ctx, transaction)
	if err != nil {
		transactionStatus = model.TransactionStatusFailed
	} else {
		transactionStatus = model.TransactionStatusPending
	}

	transaction.Status = transactionStatus
	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	err = i.rabbitMQ.PublishWithContext(ctx,
		usecase.TransactionExchangeName,
		usecase.TransactionReservedRoutingKey,
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

func (i UseCase) reserveUserBalance(ctx context.Context, transaction model.Transaction) error {
	totalPrice, err := transaction.GetTotalPrice()
	if err != nil {
		return fmt.Errorf("get total price: %w", err)
	}

	err = i.userRepository.DeductUserBalance(transaction.UserId, totalPrice)
	if err != nil {
		return fmt.Errorf("deduct user balance: %w", err)
	}

	return nil
}
