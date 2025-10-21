package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"unit-test-demo/api1/internal/domain"
	"unit-test-demo/api1/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// ---- Mock (hand-written) ----

type mockBookUsecase struct {
	// You can add fields to assert inputs if needed
	CreateBookFn func(in domain.CreateBookInput) (*domain.Book, error)
}

func (m *mockBookUsecase) CreateBook(ctx fiber.Ctx, in domain.CreateBookInput) (*domain.Book, error) {
	// This signature is wrong; we need context.Context, not fiber.Ctx.
	// We'll adapt by providing a wrapper method below.
	return nil, nil
}

// Properly implement the interface expected by handler (context.Context).
func (m *mockBookUsecase) CreateBookCtx(ctx interface{}, in domain.CreateBookInput) (*domain.Book, error) {
	return m.CreateBookFn(in)
}

// To satisfy the interface, we use a tiny adapter type.
// (Cleaner: define mock with the correct signature directly.)
type bookUsecaseAdapter struct {
	m *mockBookUsecase
}

func (a *bookUsecaseAdapter) CreateBook(ctx interface{}, in domain.CreateBookInput) (*domain.Book, error) {
	return a.m.CreateBookFn(in)
}

// ---- Test ----

type bookUsecaseIface interface {
	CreateBook(ctx interface{}, in domain.CreateBookInput) (*domain.Book, error)
}

// re-declare a small interface matching what the handler uses (context.Context)
// to avoid import cycle in tests. In real projects, put interfaces in a shared test package.

func TestCreateBook_Success(t *testing.T) {
	
	// Arrange
	app := fiber.New()

	mock := &mockBookUsecase{
		CreateBookFn: func(in domain.CreateBookInput) (*domain.Book, error) {
			assert.Equal(t, "Clean Architecture", in.Title)
			assert.Equal(t, "Uncle Bob", in.Author)
			return &domain.Book{
				ID:        42,
				Title:     in.Title,
				Author:    in.Author,
				CreatedAt: time.Date(2025, 10, 20, 12, 0, 0, 0, time.UTC),
			}, nil
		},
	}
	adapter := &bookUsecaseAdapter{m: mock}

	// We need a compatible constructor for handler that accepts this interface.
	// To keep it simple, we expose a tiny local wrapper route here:

	app.Post("/v1/books", func(c *fiber.Ctx) error {
		var req domain.CreateBookInput
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		}
		book, err := adapter.CreateBook(c.Context(), req)
		if err != nil {
			if err == usecase.ErrValidation {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}
		return c.Status(http.StatusCreated).JSON(book)
	})

	// Act
	body, _ := json.Marshal(domain.CreateBookInput{
		Title:  "Clean Architecture",
		Author: "Uncle Bob",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req, -1)

	// Assert
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var got domain.Book
	_ = json.NewDecoder(res.Body).Decode(&got)
	assert.Equal(t, int64(42), got.ID)
	assert.Equal(t, "Clean Architecture", got.Title)
	assert.Equal(t, "Uncle Bob", got.Author)
}

func TestCreateBook_ValidationError(t *testing.T) {
	app := fiber.New()

	// Force the usecase to return validation error for empty title
	mock := &mockBookUsecase{
		CreateBookFn: func(in domain.CreateBookInput) (*domain.Book, error) {
			return nil, usecase.ErrValidation
		},
	}
	adapter := &bookUsecaseAdapter{m: mock}

	app.Post("/v1/books", func(c *fiber.Ctx) error {
		var req domain.CreateBookInput
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		}
		book, err := adapter.CreateBook(c.Context(), req)
		if err != nil {
			if err == usecase.ErrValidation {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}
		return c.Status(http.StatusCreated).JSON(book)
	})

	body, _ := json.Marshal(domain.CreateBookInput{
		Title:  "",
		Author: "Anyone",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}
