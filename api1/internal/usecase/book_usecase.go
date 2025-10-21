package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"unit-test-demo/api1/internal/domain"
)

var (
	ErrValidation = errors.New("validation error")
)

type BookUsecase interface {
	CreateBook(ctx context.Context, in domain.CreateBookInput) (*domain.Book, error)
}

type bookUsecase struct {
	repo domain.BookRepository
}

func NewBookUsecase(repo domain.BookRepository) BookUsecase {
	return &bookUsecase{repo: repo}
}

func (u *bookUsecase) CreateBook(ctx context.Context, in domain.CreateBookInput) (*domain.Book, error) {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Author) == "" {
		return nil, ErrValidation
	}

	book, err := u.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	if book.CreatedAt.IsZero() {
		book.CreatedAt = time.Now()
	}

	return book, nil
}
