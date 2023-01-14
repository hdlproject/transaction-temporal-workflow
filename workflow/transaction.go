package workflow

import (
	"fmt"
	"time"

	"transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/model"

	"go.temporal.io/sdk/workflow"
)

type (
	Transaction struct {
		transactionActivity activity.Transaction
	}
)

func NewTransaction(transactionActivity activity.Transaction) Transaction {
	return Transaction{
		transactionActivity: transactionActivity,
	}
}

func (i Transaction) CreateTransaction(ctx workflow.Context, transaction model.Transaction) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.transactionActivity.CreateTransaction, transaction).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("execute activity: %w", err)
	}
	return nil
}

func (i Transaction) ProcessTransaction(ctx workflow.Context, transactionId string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.transactionActivity.ProcessTransaction, transactionId).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("execute activity: %w", err)
	}
	return nil
}
