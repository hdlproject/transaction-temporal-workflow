package cmd

import (
	amqp "github.com/rabbitmq/amqp091-go"

	internalActivity "transaction-temporal-workflow/activity"
	"transaction-temporal-workflow/dependency"
	internalRepository "transaction-temporal-workflow/repository"
	"transaction-temporal-workflow/usecase/idempotency"
	"transaction-temporal-workflow/usecase/transaction"
	"transaction-temporal-workflow/usecase/user"
	internalWorkflow "transaction-temporal-workflow/workflow"
)

var (
	TransactionActivity internalActivity.Transaction
	UserActivity        internalActivity.User
)

var (
	TransactionWorkflow internalWorkflow.Transaction
	UserWorkflow        internalWorkflow.User
)

var (
	TransactionUseCase transaction.UseCase
	UserUseCase        user.UseCase
	IdempotencyUseCase idempotency.UseCase
)

var (
	TransactionRepository internalRepository.Transaction
	UserRepository        internalRepository.User
	IdempotencyRepository internalRepository.Idempotency

	RabbitMQ *amqp.Channel
)

var (
	TransactionExchangeName      = "transaction"
	TransactionCreatedRoutingKey = "transaction.created"

	UserServiceQueueName = "user_service"
)

func init() {
	redis := dependency.NewRedis()
	db := dependency.NewPostgreSQL()
	RabbitMQ = dependency.NewRabbitMQ()
	initRabbitMQ(RabbitMQ)

	TransactionRepository = internalRepository.NewTransaction(db, RabbitMQ)
	UserRepository = internalRepository.NewUser(db)
	IdempotencyRepository = internalRepository.NewIdempotency(redis)

	TransactionUseCase = transaction.NewUseCase(
		TransactionRepository,
		IdempotencyRepository,
		RabbitMQ,
	)
	UserUseCase = user.NewUseCase(
		IdempotencyUseCase,
		UserRepository,
		RabbitMQ,
	)
	IdempotencyUseCase = idempotency.NewUseCase(
		IdempotencyRepository,
	)

	TransactionActivity = internalActivity.NewTransaction(
		TransactionUseCase,
	)
	UserActivity = internalActivity.NewUser(
		UserUseCase,
	)

	TransactionWorkflow = internalWorkflow.NewTransaction(
		TransactionActivity,
	)
	UserWorkflow = internalWorkflow.NewUser(
		UserActivity,
	)
}

func initRabbitMQ(rabbitMQ *amqp.Channel) {
	dependency.AddExchange(rabbitMQ, TransactionExchangeName)
	dependency.AddQueue(rabbitMQ, UserServiceQueueName)
	dependency.AddRouting(rabbitMQ, TransactionExchangeName, UserServiceQueueName, TransactionCreatedRoutingKey)
}
