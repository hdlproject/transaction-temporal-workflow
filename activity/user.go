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

func (i User) ReserveUserBalance(ctx context.Context, transaction model.Transaction) error {
	err := i.userUseCase.ReserveUserBalance(ctx, transaction)
	if err != nil {
		return fmt.Errorf("reserve user balance: %w", err)
	}

	return nil
}

func (i User) PublishUserBalanceEvent(ctx context.Context) error {
	err := i.userUseCase.PublishUserBalanceEvent(ctx)
	if err != nil {
		return fmt.Errorf("publish user balance event: %w", err)
	}

	return nil
}
