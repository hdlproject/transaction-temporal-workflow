package user

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/repository"
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

func (i UseCase) ProcessTransaction(ctx context.Context, transaction model.Transaction, idempotencyKey string) error {
	isAllowed, err := i.idempotencyUseCase.IsAllowed(idempotencyKey)
	if err != nil {
		return fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil
	}

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
