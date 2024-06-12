package application

import (
	"library-api/internal/book/domain"
)

type BookRepository interface {
	FindAll() ([]domain.Book, error)
	FindByID(id *domain.Id) (*domain.Book, error)
	Create(book *domain.Book) error
	Update(book *domain.Book) error
}
