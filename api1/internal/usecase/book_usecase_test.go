package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"unit-test-demo/api1/internal/domain"
	domain_mock "unit-test-demo/api1/internal/mocks/domain"
	"unit-test-demo/api1/internal/usecase"
)

func TestCreateBook_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := domain_mock.NewMockBookRepository(ctrl)
	in := domain.CreateBookInput{Title: "X", Author: "Y"}
	out := &domain.Book{ID: 99, Title: "X", Author: "Y", CreatedAt: time.Now()}

	mr.EXPECT().Create(gomock.Any(), in).Return(out, nil)

	uc := usecase.NewBookUsecase(mr)
	got, err := uc.CreateBook(context.Background(), in)

	assert.NoError(t, err)
	assert.Equal(t, int64(99), got.ID)
}

func TestCreateBook_Validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := domain_mock.NewMockBookRepository(ctrl)
	uc := usecase.NewBookUsecase(mr)

	got, err := uc.CreateBook(context.Background(), domain.CreateBookInput{Title: "", Author: "Y"})
	assert.ErrorIs(t, err, usecase.ErrValidation)
	assert.Nil(t, got)
}

func TestCreateBook_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := domain_mock.NewMockBookRepository(ctrl)
	in := domain.CreateBookInput{Title: "A", Author: "B"}

	mr.EXPECT().Create(gomock.Any(), in).Return((*domain.Book)(nil), errors.New("db down"))

	uc := usecase.NewBookUsecase(mr)
	got, err := uc.CreateBook(context.Background(), in)

	assert.Error(t, err)
	assert.Nil(t, got)
}
