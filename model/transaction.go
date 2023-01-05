package model

type (
	TransactionStatus string

	Transaction struct {
		Id     string            `json:"id"`
		Status TransactionStatus `json:"status"`
	}
)

const TransactionStatusCreated TransactionStatus = "CREATED"
const TransactionStatusPending TransactionStatus = "PENDING"
const TransactionStatusSuccess TransactionStatus = "SUCCESS"
const TransactionStatusFailed TransactionStatus = "FAILED"
