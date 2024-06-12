package memory

import (
	"errors"
	"library-api/internal/book/application"
	"library-api/internal/book/domain"
)

type memoryBookRepository struct {
	books map[domain.Id]domain.Book
}

func NewInMemoryBookRepository() application.BookRepository {
	books := make(map[domain.Id]domain.Book)

	return &memoryBookRepository{
		books: books,
	}
}

func (r *memoryBookRepository) FindAll() ([]domain.Book, error) {
	list := make([]domain.Book, 0, len(r.books))
	for _, book := range r.books {
		list = append(list, book)
	}
	return list, nil
}

func (r *memoryBookRepository) FindByID(id *domain.Id) (*domain.Book, error) {
	book, exists := r.books[*id]
	if !exists {
		return nil, nil
	}
	return &book, nil
}

func (r *memoryBookRepository) Create(book *domain.Book) error {
	_, exists := r.books[*book.ID]
	if exists {
		return errors.New("book already exists")
	}

	r.books[*book.ID] = *book
	return nil
}

func (r *memoryBookRepository) Update(book *domain.Book) error {
	r.books[*book.ID] = *book
	return nil
}
