package user

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"temporalio-poc/model"
	"temporalio-poc/repository"
	"temporalio-poc/usecase"
	"temporalio-poc/usecase/idempotency"
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
	err := i.userRepository.DeductUserBalance(transaction)
	if err != nil {
		return fmt.Errorf("deduct user balance: %w", err)
	}

	return nil
}

func (i UseCase) PublishUserBalanceEvent(ctx context.Context) error {
	userBalanceEvents, err := i.userRepository.GetUnpublishedUserBalanceEvents()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("get unpublished user balance events: %w", err)
		}

		return nil
	}

	for _, userBalanceEvent := range userBalanceEvents {
		userBalanceEventJson, err := json.Marshal(userBalanceEvent)
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
				Body:        userBalanceEventJson,
			},
		)
		if err != nil {
			return fmt.Errorf("rabbitmq publish with context: %w", err)
		}

		err = i.userRepository.PublishUserBalanceEvent(userBalanceEvent.Id)
		if err != nil {
			return fmt.Errorf("publish transaction: %w", err)
		}
	}

	return nil
}
