package model

type User struct {
	Id      string `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (User) TableName() string {
	return "user"
}
