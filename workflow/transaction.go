package workflow

import (
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

func (i Transaction) CreateTransaction(ctx workflow.Context, transaction model.Transaction, idempotencyKey string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.transactionActivity.CreateTransaction, transaction, idempotencyKey).Get(ctx, nil)
	return err
}

func (i Transaction) ProcessTransaction(ctx workflow.Context, transactionId, idempotencyKey string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, i.transactionActivity.ProcessTransaction, transactionId, idempotencyKey).Get(ctx, nil)
	return err
}
