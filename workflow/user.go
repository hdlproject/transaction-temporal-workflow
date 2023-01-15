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

func (i User) ReserveUserBalance(ctx workflow.Context, transaction model.Transaction) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.userActivity.ReserveUserBalance, transaction).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("execute activity: %w", err)
	}

	return nil
}
