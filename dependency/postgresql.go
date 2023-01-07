package dependency

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQL() *gorm.DB {
	// TODO: move credential to .env
	dsn := `host=localhost 
			user=app 
			password=app 
			dbname=app 
			port=5433 
			sslmode=disable 
			TimeZone=Asia/Jakarta`

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("gorm open: " + err.Error())
	}

	return db
}
