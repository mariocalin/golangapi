package book

import (
	"errors"

	"github.com/google/uuid"
)

type BookCommand struct {
	Name        *Name
	PublishDate *PublishDate
	Categories  *Categories
	Description *Description
}

func (command *BookCommand) isValidForCreate() bool {
	return command.Name != nil && command.PublishDate != nil && command.Categories != nil && command.Description != nil
}

type BookService interface {
	GetBooks() ([]Book, error)
	GetBookByID(id *BookId) (*Book, error)
	CreateBook(command *BookCommand) (*Book, error)
	UpdateBook(id *BookId, command *BookCommand) error
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

func (s *service) GetBookByID(id *BookId) (*Book, error) {
	return s.repo.FindByID(id)
}

func (s *service) CreateBook(command *BookCommand) (*Book, error) {

	if !command.isValidForCreate() {
		return nil, errors.New("cannot create book")
	}

	id := uuid.New()

	newBook := Book{
		ID:          &id,
		Name:        command.Name,
		PublishDate: command.PublishDate,
		Categories:  command.Categories,
		Description: command.Description,
	}

	s.repo.Create(&newBook)
	s.events.PropagateBookCreated(&BookCreated{Id: newBook.ID, Name: newBook.Name})

	return &newBook, nil
}

func (s *service) UpdateBook(id *BookId, command *BookCommand) error {
	existingBook, notFound := s.repo.FindByID(id)
	if notFound != nil {
		return notFound
	}

	if command.Name != nil {
		existingBook.Name = command.Name
	}

	if command.PublishDate != nil {
		existingBook.PublishDate = command.PublishDate
	}

	if command.Categories != nil {
		existingBook.Categories = command.Categories
	}

	if command.Description != nil {
		existingBook.Description = command.Description
	}

	s.repo.Update(existingBook)

	return nil
}
