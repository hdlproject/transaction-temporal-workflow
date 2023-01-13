package cmd

import (
	amqp "github.com/rabbitmq/amqp091-go"

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

var (
	TransactionExchangeName      = "transaction"
	TransactionCreatedRoutingKey = "transaction.created"

	UserServiceQueueName = "user_service"
)

func init() {
	redis := dependency.NewRedis()
	db := dependency.NewPostgreSQL()
	rabbitMQ := dependency.NewRabbitMQ()
	initRabbitMQ(rabbitMQ)

	TransactionRepository = internalRepository.NewTransaction(db, rabbitMQ)
	IdempotencyRepository = internalRepository.NewIdempotency(redis)

	TransactionUseCase = transaction.NewUseCase(
		TransactionRepository,
		IdempotencyRepository,
		rabbitMQ,
	)

	TransactionActivity = internalActivity.NewTransaction(
		TransactionUseCase,
	)

	TransactionWorkflow = internalWorkflow.NewTransaction(
		TransactionActivity,
	)
}

func initRabbitMQ(rabbitMQ *amqp.Channel) {
	dependency.AddExchange(rabbitMQ, TransactionExchangeName)
	dependency.AddQueue(rabbitMQ, UserServiceQueueName)
	dependency.AddRouting(rabbitMQ, TransactionExchangeName, UserServiceQueueName, TransactionCreatedRoutingKey)
}
