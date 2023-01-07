package cmd

import (
	internalActivity "transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/dependency"
	internalRepository "transaction-temporal-workflow/repository"
	internalWorkflow "transaction-temporal-workflow/workflow"
)

var (
	TransactionActivity internalActivity.Transaction
)

var (
	TransactionWorkflow internalWorkflow.Transaction
)

func init() {
	redis := dependency.NewRedis()
	db := dependency.NewPostgreSQL()

	TransactionActivity = internalActivity.NewTransaction(
		internalRepository.NewTransaction(
			redis,
			db,
		),
	)

	TransactionWorkflow = internalWorkflow.NewTransaction(
		TransactionActivity,
	)
}
