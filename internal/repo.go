package internal

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("entity not found")

//go:generate mockery --name=BookRepository --output=mocks --filename=book_repository_mock.go --with-expecter=true
type BookRepository interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindByID(ctx context.Context, id BookId) (Book, error)
	Create(ctx context.Context, book Book) error
	Update(ctx context.Context, book Book) error
}
