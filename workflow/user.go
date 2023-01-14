package workflow

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.temporal.io/sdk/workflow"

	"transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/cmd"
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

func (i User) ConsumeMessage(ctx workflow.Context, message amqp091.Delivery) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	switch message.RoutingKey {
	case cmd.TransactionCreatedRoutingKey:
		var transaction model.Transaction
		err := json.Unmarshal(message.Body, &transaction)
		if err != nil {
			return fmt.Errorf("json unmarshal: %w", err)
		}

		err = workflow.ExecuteActivity(ctx, i.userActivity.ProcessTransaction, transaction, transaction.TransactionId).Get(ctx, nil)
		if err != nil {
			return fmt.Errorf("execute activity: %w", err)
		}
	}

	return nil
}
