package model

import (
	"fmt"
	"time"
)

type (
	TransactionStatus string

	Transaction struct {
		Id            int64             `json:"id" gorm:"primaryKey"`
		TransactionId string            `json:"transaction_id"`
		Status        TransactionStatus `json:"status"`
		Amount        int               `json:"amount"`
		ProductCode   string            `json:"product_code"`
		Product       Product           `gorm:"foreignKey:ProductCode"`
		UserId        string            `json:"user_id" `
		User          User              `gorm:"foreignKey:UserId"`
		CreatedAt     time.Time         `json:"created_at"`
		IsPublished   bool              `json:"is_published"`
	}

	TransactionQuery struct {
		Transaction
	}
)

const TransactionStatusCreated TransactionStatus = "CREATED"
const TransactionStatusPending TransactionStatus = "PENDING"
const TransactionStatusSuccess TransactionStatus = "SUCCESS"
const TransactionStatusFailed TransactionStatus = "FAILED"

func (t *Transaction) TableName() string {
	return "transaction"
}

func (t *TransactionQuery) TableName() string {
	return "transaction_query"
}

func (t *Transaction) GetTotalPrice() (int, error) {
	if t.Product.Code == "" {
		return 0, fmt.Errorf("product is empty")
	}
	return t.Product.Price * t.Amount, nil
}
