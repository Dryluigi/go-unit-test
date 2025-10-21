package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"unit-test-demo/api1/internal/domain"
	"unit-test-demo/api1/internal/usecase"
)

// ---- Fake repository (in-memory) ----

type fakeRepo struct {
	// seed & state
	created []*domain.Book

	// knobs
	nextErr error
	// if set, repo will return CreatedAt zero to test usecase autofill
	returnZeroCreatedAt bool

	// observability
	createCalls int
	lastInput   domain.CreateBookInput
}

func newFakeRepo() *fakeRepo { return &fakeRepo{} }

func (f *fakeRepo) Create(ctx context.Context, in domain.CreateBookInput) (*domain.Book, error) {
	f.createCalls++
	f.lastInput = in

	if f.nextErr != nil {
		err := f.nextErr
		f.nextErr = nil
		return nil, err
	}

	var createdAt time.Time
	if !f.returnZeroCreatedAt {
		createdAt = time.Now()
	}

	b := &domain.Book{
		ID:        int64(len(f.created) + 1),
		Title:     in.Title,
		Author:    in.Author,
		CreatedAt: createdAt,
	}
	f.created = append(f.created, b)
	return b, nil
}

// ---- Tests ----

func TestBookUsecase_CreateBook_HappyPath(t *testing.T) {
	// Arrange
	fake := newFakeRepo()
	uc := usecase.NewBookUsecase(fake)
	in := domain.CreateBookInput{
		Title:  "Clean Architecture",
		Author: "Robert C. Martin",
	}

	// Act
	got, err := uc.CreateBook(context.Background(), in)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("got nil book")
	}
	if got.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}
	if got.Title != in.Title || got.Author != in.Author {
		t.Fatalf("book fields mismatch: got=%+v want title=%q author=%q", got, in.Title, in.Author)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set, got zero")
	}
	if fake.createCalls != 1 {
		t.Fatalf("expected 1 repo call, got %d", fake.createCalls)
	}
}

func TestBookUsecase_CreateBook_ValidationError(t *testing.T) {
	// Arrange
	fake := newFakeRepo()
	uc := usecase.NewBookUsecase(fake)

	// Act
	_, err := uc.CreateBook(context.Background(), domain.CreateBookInput{
		Title:  "   ", // invalid
		Author: "Author",
	})

	// Assert
	if !errors.Is(err, usecase.ErrValidation) {
		t.Fatalf("want ErrValidation, got %v", err)
	}
	if fake.createCalls != 0 {
		t.Fatalf("repo should not be called on validation error, calls=%d", fake.createCalls)
	}
}

func TestBookUsecase_CreateBook_RepoError(t *testing.T) {
	// Arrange
	fake := newFakeRepo()
	fake.nextErr = errors.New("db timeout")
	uc := usecase.NewBookUsecase(fake)

	// Act
	_, err := uc.CreateBook(context.Background(), domain.CreateBookInput{
		Title:  "Dune",
		Author: "Frank Herbert",
	})

	// Assert
	if err == nil || err.Error() != "db timeout" {
		t.Fatalf("expected db timeout error, got %v", err)
	}
	if fake.createCalls != 1 {
		t.Fatalf("expected exactly 1 repo call, got %d", fake.createCalls)
	}
}

func TestBookUsecase_CreateBook_FillCreatedAt_WhenRepoReturnsZero(t *testing.T) {
	// Arrange
	fake := newFakeRepo()
	fake.returnZeroCreatedAt = true // simulate repo not setting CreatedAt
	uc := usecase.NewBookUsecase(fake)

	// Act
	start := time.Now()
	got, err := uc.CreateBook(context.Background(), domain.CreateBookInput{
		Title:  "The Pragmatic Programmer",
		Author: "Andrew Hunt",
	})
	end := time.Now()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be filled by usecase when repo returns zero")
	}
	// sanity window check (±2s)
	if got.CreatedAt.Before(start.Add(-2*time.Second)) || got.CreatedAt.After(end.Add(2*time.Second)) {
		t.Fatalf("CreatedAt not within expected window, got=%v start=%v end=%v", got.CreatedAt, start, end)
	}
}
