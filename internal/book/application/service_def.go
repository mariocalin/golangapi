package application

import (
	"library-api/internal/book/domain"
)

//go:generate mockery --name=Service --outpkg=book --dir=. --filename=service_mock.go --output=. --structname=MockService
type Service interface {
	GetBooks() ([]domain.Book, error)
	GetBookByID(id *domain.Id) (*domain.Book, error)
	CreateBook(command *CreateCommand) (*domain.Book, error)
	UpdateBook(id *domain.Id, command *CreateCommand) error
}

type CreateCommand struct {
	Name        *domain.Name
	PublishDate *domain.PublishDate
	Categories  *domain.Categories
	Description *domain.Description
}

func (command *CreateCommand) isValidForCreate() bool {
	return command.Name != nil && command.PublishDate != nil && command.Categories != nil && command.Description != nil
}
