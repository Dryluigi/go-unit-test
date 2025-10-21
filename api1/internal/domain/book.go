package domain

import (
	"context"
	"time"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookRepository interface {
	Create(ctx context.Context, in CreateBookInput) (*Book, error)
}
