package main

type PostPayload struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Status     string `json:"status" validate:"required"`
	IsFeatured bool   `json:"is_featured"`
	// Tags    []models.Tag `json:"tags"` // TODO: input string: misal python, golang, dst.
}

type Validation struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}
