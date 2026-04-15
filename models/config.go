package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	GetDB()
}

// GetDB function to get database connection
func GetDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Tag{}, &Post{})
}
