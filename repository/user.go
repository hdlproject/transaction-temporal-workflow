package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return User{
		db: db,
	}
}

func (i User) DeductUserBalance(userId string, amount int) error {
	result := i.db.Exec("UPDATE user SET balance = ? WHERE id = ?", gorm.Expr("balance - ?", amount), userId)
	if result.Error != nil {
		return fmt.Errorf("deduct user balance: %w", result.Error)
	}

	return nil
}
