package models

import "gorm.io/gorm"

type User struct {
	ID        int64   `gorm:"primaryKey" json:"id"`
	FirstName string  `gorm:"size:20" json:"first_name"`
	LastName  string  `gorm:"size:20" json:"last_name"`
	Username  string  `gorm:"size:20;unique" json:"username"`
	Email     string  `gorm:"size:50;unique" json:"email"`
	Photo     *string `gorm:"size:255" json:"photo,omitempty"`
	Password  string  `gorm:"size:128" json:"password"`
}

type UserModel interface {
	// CreateUser function to create new user.
	CreateUser(user *User) error
	// GetUserByID function to get user by id
	GetUserByID(id int64) (User, error)
	// GetUserByUsername function to get user by username.
	GetUserByUsername(username string) (User, error)
}

type userModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) UserModel {
	return userModel{db: db}
}

func (u userModel) CreateUser(user *User) error {
	return u.db.Model(&User{}).Create(user).Error
}

func (u userModel) GetUserByID(id int64) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("id = ?", id).First(&user).Error
	return user, err
}

func (u userModel) GetUserByUsername(username string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("username = ?", username).First(&user).Error
	return user, err
}
