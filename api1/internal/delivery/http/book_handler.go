package http

import (
	"net/http"

	"unit-test-demo/api1/internal/domain"
	"unit-test-demo/api1/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	uc usecase.BookUsecase
}

func NewBookHandler(r fiber.Router, uc usecase.BookUsecase) {
	h := &BookHandler{uc: uc}
	v1 := r.Group("/v1")
	v1.Post("/books", h.CreateBook)
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req domain.CreateBookInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body",
		})
	}

	book, err := h.uc.CreateBook(c.Context(), req)
	if err != nil {
		if err == usecase.ErrValidation {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	return c.Status(http.StatusCreated).JSON(book)
}
