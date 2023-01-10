package cmd

import (
	internalActivity "transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/dependency"
	internalRepository "transaction-temporal-workflow/repository"
	"transaction-temporal-workflow/usecase/transaction"
	internalWorkflow "transaction-temporal-workflow/workflow"
)

var (
	TransactionActivity internalActivity.Transaction
)

var (
	TransactionWorkflow internalWorkflow.Transaction
)

var (
	TransactionUseCase transaction.UseCase
)

func init() {
	redis := dependency.NewRedis()
	db := dependency.NewPostgreSQL()

	TransactionUseCase = transaction.NewUseCase(
		internalRepository.NewTransaction(
			redis,
			db,
		),
	)

	TransactionActivity = internalActivity.NewTransaction(
		TransactionUseCase,
	)

	TransactionWorkflow = internalWorkflow.NewTransaction(
		TransactionActivity,
	)
}
