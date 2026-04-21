package main

type PostPayload struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Status     string `json:"status" validate:"required"`
	IsFeatured bool   `json:"is_featured"`
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
