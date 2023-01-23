package model

import "time"

type User struct {
	Id      string `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type UserBalanceEvent struct {
	Id                int64             `json:"id" gorm:"primaryKey"`
	UserId            string            `json:"user_id"`
	User              User              `gorm:"foreignKey:UserId"`
	Balance           int               `json:"balance"`
	TransactionId     string            `json:"transaction_id"`
	TransactionStatus TransactionStatus `json:"transaction_status"`
	CreatedAt         time.Time         `json:"created_at"`
	IsPublished       bool              `json:"is_published"`
}

func (User) TableName() string {
	return "user"
}

func (UserBalanceEvent) TableName() string {
	return "user_balance_event"
}
