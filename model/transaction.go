package model

import "time"

type (
	TransactionStatus string

	Transaction struct {
		Id          string            `json:"id"`
		Status      TransactionStatus `json:"status"`
		Amount      int               `json:"amount"`
		ProductCode string            `json:"product_code"`
		UserId      string            `json:"user_id"`
		CreatedAt   time.Time         `json:"created_at"`
	}
)

const TransactionStatusCreated TransactionStatus = "CREATED"
const TransactionStatusPending TransactionStatus = "PENDING"
const TransactionStatusSuccess TransactionStatus = "SUCCESS"
const TransactionStatusFailed TransactionStatus = "FAILED"
