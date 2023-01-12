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

var (
	TransactionRepository internalRepository.Transaction
	IdempotencyRepository internalRepository.Idempotency
)

func init() {
	redis := dependency.NewRedis()
	db := dependency.NewPostgreSQL()

	TransactionRepository = internalRepository.NewTransaction(db)
	IdempotencyRepository = internalRepository.NewIdempotency(redis)

	TransactionUseCase = transaction.NewUseCase(
		TransactionRepository,
		IdempotencyRepository,
	)

	TransactionActivity = internalActivity.NewTransaction(
		TransactionUseCase,
	)

	TransactionWorkflow = internalWorkflow.NewTransaction(
		TransactionActivity,
	)
}
