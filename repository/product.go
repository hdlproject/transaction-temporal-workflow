package repository

import (
	"fmt"

	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) Product {
	return Product{
		db: db,
	}
}

func (i Product) GetProductByCode(code string) (product model.Product, err error) {
	result := i.db.First(&product, "code = ?", code)
	if result.Error != nil {
		return model.Product{}, fmt.Errorf("get product by code: %w", result.Error)
	}

	return product, nil
}
