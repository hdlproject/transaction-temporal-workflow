package model

type Product struct {
	Code  string `json:"code" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (Product) TableName() string {
	return "product"
}
