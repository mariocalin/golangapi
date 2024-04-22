package book

import (
	"errors"
)

type BookRepository interface {
	FindAll() ([]Book, error)
	FindByID(id BookId) (*Book, error)
	Create(book *Book) error
}

type memoryBookRepository struct {
	books map[BookId]Book
}

func NewRepository() BookRepository {
	books := make(map[BookId]Book)

	return &memoryBookRepository{
		books: books,
	}
}

func (r *memoryBookRepository) FindAll() ([]Book, error) {
	list := make([]Book, 0, len(r.books))
	for _, book := range r.books {
		list = append(list, book)
	}
	return list, nil
}

func (r *memoryBookRepository) FindByID(id BookId) (*Book, error) {
	book, exists := r.books[id]
	if !exists {
		return nil, nil // or your error
	}
	return &book, nil
}

func (r *memoryBookRepository) Create(book *Book) error {
	// Check if the book already exists
	_, exists := r.books[book.ID]
	if exists {
		return errors.New("book already exists")
	}

	// If the book doesn't exist, create it
	r.books[book.ID] = *book
	return nil
}
