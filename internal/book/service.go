package book

import "github.com/google/uuid"

type CreateBookCommand struct {
	Name        Name
	PublishDate PublishDate
	Categories  Categories
	Description Description
}

type UpdateBookCommand struct {
	Name        *Name
	PublishDate *PublishDate
	Categories  *Categories
	Description *Description
}

type BookService interface {
	GetBooks() ([]Book, error)
	GetBookByID(id BookId) (*Book, error)
	CreateBook(command *CreateBookCommand) (*Book, error)
	UpdateBook(id BookId, command *UpdateBookCommand) error
}

type service struct {
	repo   BookRepository
	events BookEventPropagator
}

func NewService(repo BookRepository, events BookEventPropagator) BookService {
	return &service{
		repo:   repo,
		events: events,
	}
}

func (s *service) GetBooks() ([]Book, error) {
	return s.repo.FindAll()
}

func (s *service) GetBookByID(id BookId) (*Book, error) {
	return s.repo.FindByID(id)
}

func (s *service) CreateBook(command *CreateBookCommand) (*Book, error) {
	newBook := Book{
		ID:          uuid.New(),
		Name:        command.Name,
		PublishDate: command.PublishDate,
		Categories:  command.Categories,
		Description: command.Description,
	}

	s.repo.Create(&newBook)
	s.events.propagateBookCreated(&BookCreated{Id: newBook.ID})

	return &newBook, nil
}

func (s *service) UpdateBook(id BookId, command *UpdateBookCommand) error {
	existingBook, notFound := s.repo.FindByID(id)
	if notFound != nil {
		return notFound
	}

	if command.Name != nil {
		existingBook.Name = *command.Name
	}

	if command.PublishDate != nil {
		existingBook.PublishDate = *command.PublishDate
	}

	if command.Categories != nil {
		existingBook.Categories = *command.Categories
	}

	if command.Description != nil {
		existingBook.Description = *command.Description
	}

	s.repo.Update(existingBook)

	return nil
}
