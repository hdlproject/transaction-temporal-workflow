package activity

import (
	"context"
	"fmt"

	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/usecase/user"
)

type (
	User struct {
		userUseCase user.UseCase
	}
)

func NewUser(userUseCase user.UseCase) User {
	return User{
		userUseCase: userUseCase,
	}
}

func (i User) ProcessTransaction(ctx context.Context, transaction model.Transaction, idempotencyKey string) error {
	err := i.userUseCase.ProcessTransaction(ctx, transaction, idempotencyKey)
	if err != nil {
		return fmt.Errorf("process transaction: %w", err)
	}

	return nil
}
