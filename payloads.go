package main

import "mime/multipart"

type PostPayload struct {
	Title      string                `form:"title" validate:"required"`
	Logo       *multipart.FileHeader `form:"logo" validate:"required"`
	Content    string                `form:"content" validate:"required"`
	Status     string                `form:"status" validate:"required"`
	IsFeatured bool                  `form:"is_featured"`
	// Tags    []models.Tag `json:"tags"` // TODO: input string: misal python, golang, dst.
}

type PostEditPayload struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Status     *string `json:"status"`
	IsFeatured *bool   `json:"is_featured"`
}

type Validation struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}
