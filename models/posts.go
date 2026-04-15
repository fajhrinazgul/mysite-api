package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"size:20"`
}

type Post struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:100" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    string    `gorm:"size:9;default:drafted" json:"status"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Tags      []Tag     `gorm:"many2many:post_tags" json:"tags"`
}

type PostModel interface {
	CreatePost(post *Post) error
}

type postModel struct{ db *gorm.DB }

func NewPostModel(db *gorm.DB) PostModel { return postModel{db: db} }

func (p postModel) CreatePost(post *Post) error {
	return p.db.Model(&Post{}).Create(post).Error
}
