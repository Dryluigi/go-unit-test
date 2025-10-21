package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"unit-test-demo/api1/internal/domain"
	"unit-test-demo/api1/internal/usecase"
)

// ---- Stub repository ----
// Returns pre-configured values, no state/logic.
type stubRepo struct {
	// canned outputs
	returnBook *domain.Book
	returnErr  error

	// optional spy bit to observe if Create was invoked
	called bool
	lastIn domain.CreateBookInput
}

func (s *stubRepo) Create(ctx context.Context, in domain.CreateBookInput) (*domain.Book, error) {
	s.called = true
	s.lastIn = in
	return s.returnBook, s.returnErr
}

// ---- Tests ----

func TestBookUsecase_CreateBook_HappyPath_WithStub(t *testing.T) {
	// Arrange: stub will succeed and return a fixed book
	now := time.Date(2025, 10, 21, 12, 0, 0, 0, time.UTC)
	stub := &stubRepo{
		returnBook: &domain.Book{
			ID:        42,
			Title:     "Clean Code",
			Author:    "Robert C. Martin",
			CreatedAt: now, // already set by stub (pretend DB populated it)
		},
	}
	uc := usecase.NewBookUsecase(stub)
	in := domain.CreateBookInput{
		Title:  "Clean Code",
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
	if got.ID != 42 {
		t.Fatalf("want ID=42 got=%d", got.ID)
	}
	if got.Title != in.Title || got.Author != in.Author {
		t.Fatalf("fields mismatch: got=%+v want title=%q author=%q", got, in.Title, in.Author)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set by stub/DB")
	}
}

func TestBookUsecase_CreateBook_ValidationError_WithStub(t *testing.T) {
	// Arrange: stub would succeed if called—but UC should fail fast
	stub := &stubRepo{
		returnBook: &domain.Book{ID: 1, Title: "X", Author: "Y", CreatedAt: time.Now()},
	}
	uc := usecase.NewBookUsecase(stub)

	// Act
	_, err := uc.CreateBook(context.Background(), domain.CreateBookInput{
		Title:  "   ", // invalid
		Author: "Y",
	})

	// Assert
	if !errors.Is(err, usecase.ErrValidation) {
		t.Fatalf("want ErrValidation got=%v", err)
	}
	// Since it's a stub, we *can* also observe if UC called it (optional spy bit)
	if stub.called {
		t.Fatalf("repo should not be called on validation failure")
	}
}

func TestBookUsecase_CreateBook_RepoError_WithStub(t *testing.T) {
	// Arrange: stub returns a canned error
	stub := &stubRepo{
		returnErr: errors.New("db timeout"),
	}
	uc := usecase.NewBookUsecase(stub)

	// Act
	_, err := uc.CreateBook(context.Background(), domain.CreateBookInput{
		Title:  "Dune",
		Author: "Frank Herbert",
	})

	// Assert
	if err == nil || err.Error() != "db timeout" {
		t.Fatalf("want 'db timeout' got=%v", err)
	}
}

func TestBookUsecase_CreateBook_FillsCreatedAt_WhenStubReturnsZero(t *testing.T) {
	// Arrange: stub returns a book with zero CreatedAt to trigger UC autofill
	stub := &stubRepo{
		returnBook: &domain.Book{
			ID:        7,
			Title:     "The Pragmatic Programmer",
			Author:    "Andrew Hunt",
			CreatedAt: time.Time{}, // zero
		},
	}
	uc := usecase.NewBookUsecase(stub)

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
		t.Fatalf("expected UC to autofill CreatedAt when stub returns zero")
	}
	// sanity window check
	if got.CreatedAt.Before(start.Add(-2*time.Second)) || got.CreatedAt.After(end.Add(2*time.Second)) {
		t.Fatalf("CreatedAt not within expected window; got=%v start=%v end=%v", got.CreatedAt, start, end)
	}
}
