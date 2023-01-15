package usecase

const TransactionTaskQueue = "TRANSACTION_TASK_QUEUE"
const UserTaskQueue = "USER_TASK_QUEUE"

const (
	TransactionExchangeName = "transaction"

	TransactionCreatedRoutingKey  = "transaction.created"
	TransactionReservedRoutingKey = "transaction.reserved"

	UserServiceQueueName        = "user_service"
	TransactionServiceQueueName = "transaction_service"
)
