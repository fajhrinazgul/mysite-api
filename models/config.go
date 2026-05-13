package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	GetDB()
}

// GetDB function to get database connection
func GetDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{}, &Tag{}, &Post{})
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Tag{}, &Post{})
}
