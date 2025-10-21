package main

import (
	"context"
	"log"
	"time"

	httpdelivery "unit-test-demo/api1/internal/delivery/http"
	"unit-test-demo/api1/internal/infrastructure/postgres"
	"unit-test-demo/api1/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/book-unit-test?sslmode=disable"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer conn.Close(context.Background())

	repo := postgres.NewBookRepository(conn)
	uc := usecase.NewBookUsecase(repo)

	app := fiber.New()
	httpdelivery.NewBookHandler(app, uc)

	log.Println("listening on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
