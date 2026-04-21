package models

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Tag struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:20"`
}

type Post struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:250" json:"title"`
	Slug        string    `gorm:"size:250;uniqueIndex;not null" json:"slug"`
	Summary     string    `gorm:"size:500" json:"summary"`
	Content     string    `gorm:"type:text" json:"content"`
	IsFeatured  bool      `gorm:"default:false" json:"is_featured"`
	Status      string    `gorm:"size:9;default:drafted" json:"status"`
	ReadingTime string    `gorm:"size:20" json:"reading_time"`
	View        uint      `gorm:"default:0" json:"view"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	Tags        []Tag     `gorm:"many2many:post_tags" json:"tags"`
}

func (p *Post) BeforeSave(tx *gorm.DB) (err error) {
	generatedSlug := slug.Make(p.Title)
	finalSlug := generatedSlug

	var count int64
	counter := 1

	for {
		query := tx.Model(&Post{}).Where("slug = ?", finalSlug)
		if p.ID != 0 {
			query = query.Where("id != ?", p.ID)
		}
		query.Count(&count)

		if count == 0 {
			break
		}

		finalSlug = fmt.Sprintf("%s-%d", generatedSlug, counter)
		counter++
	}

	p.Slug = finalSlug

	// reading time
	words := strings.Fields(p.Content)
	wordCount := len(words)
	minutes := math.Ceil(float64(wordCount) / 200)
	if minutes <= 1 {
		p.ReadingTime = "1 min read"
	} else {
		p.ReadingTime = fmt.Sprintf("%.0f min read", minutes)
	}

	// auto summary
	if p.Summary == "" && p.Content != "" {
		// clean tag html (optional)
		cleanContent := p.Content
		if len(cleanContent) > 150 {
			p.Summary = cleanContent[:147] + "..."
		} else {
			p.Summary = cleanContent
		}
	}

	return nil
}

type PostModel interface {
	CreatePost(post *Post) error
	GetAllPost() []Post
	GetPagedPosts(status string, limit, offset int) ([]Post, int64, error)
	GetPostBySlug(slug string) (Post, error)
}

type postModel struct{ db *gorm.DB }

func NewPostModel(db *gorm.DB) PostModel { return postModel{db: db} }

func (p postModel) CreatePost(post *Post) error {
	return p.db.Model(&Post{}).Create(post).Error
}

func (p postModel) GetAllPost() []Post {
	var posts []Post
	p.db.Model(&Post{}).Find(&posts)
	return posts
}

func (p postModel) GetPagedPosts(status string, limit, offset int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	err := p.db.Model(&Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = p.db.Limit(limit).Offset(offset).Order("created_at DESC").Where("status = ?", status).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (p postModel) GetPostBySlug(slug string) (Post, error) {
	var post Post
	err := p.db.Model(&Post{}).Where("slug = ?", slug).First(&post).Error
	return post, err
}
