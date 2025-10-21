package postgres

import (
	"context"
	"time"

	"unit-test-demo/api1/internal/domain"

	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	conn *pgx.Conn
}

func NewBookRepository(conn *pgx.Conn) *BookRepository {
	return &BookRepository{conn: conn}
}

func (r *BookRepository) Create(ctx context.Context, in domain.CreateBookInput) (*domain.Book, error) {
	var b domain.Book
	err := r.conn.QueryRow(ctx,
		`INSERT INTO books (title, author) 
         VALUES ($1, $2)
         RETURNING id, title, author, created_at`,
		in.Title, in.Author,
	).Scan(&b.ID, &b.Title, &b.Author, &b.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Normalize timezone if needed
	b.CreatedAt = b.CreatedAt.UTC().Truncate(time.Second)
	return &b, nil
}
