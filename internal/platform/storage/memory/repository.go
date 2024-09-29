package memory

import (
	"context"
	"errors"
	"library-api/internal"
)

type memoryBookRepository struct {
	books map[internal.BookId]internal.Book
}

func NewInMemoryBookRepository() internal.BookRepository {
	books := make(map[internal.BookId]internal.Book)

	return &memoryBookRepository{
		books: books,
	}
}

func (r *memoryBookRepository) FindAll(ctx context.Context) ([]internal.Book, error) {
	list := make([]internal.Book, 0, len(r.books))
	for _, book := range r.books {
		list = append(list, book)
	}

	return list, nil
}

func (r *memoryBookRepository) FindByID(ctx context.Context, id internal.BookId) (internal.Book, error) {
	book, exists := r.books[id]
	if !exists {
		return internal.Book{}, internal.ErrNotFound
	}

	return book, nil
}

func (r *memoryBookRepository) Create(ctx context.Context, book internal.Book) error {
	_, exists := r.books[book.ID]
	if exists {
		return errors.New("book already exists")
	}

	r.books[book.ID] = book

	return nil
}

func (r *memoryBookRepository) Update(ctx context.Context, book internal.Book) error {
	_, exists := r.books[book.ID]
	if !exists {
		return errors.New("book does not exist")
	}

	r.books[book.ID] = book

	return nil
}
