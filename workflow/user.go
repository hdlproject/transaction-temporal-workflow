package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/model"
)

type (
	User struct {
		userActivity activity.User
	}
)

func NewUser(userActivity activity.User) User {
	return User{
		userActivity: userActivity,
	}
}

func (i User) ProcessTransaction(ctx workflow.Context, transaction model.Transaction, idempotencyKey string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.userActivity.ProcessTransaction, transaction, idempotencyKey).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("execute activity: %w", err)
	}

	return nil
}
