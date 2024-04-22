package book

import "github.com/google/uuid"

type CreateBookCommand struct {
	Name        Name
	PublishDate PublishDate
	Categories  Categories
	Description Description
}

type BookService interface {
	GetBooks() ([]Book, error)
	GetBookByID(id BookId) (*Book, error)
	CreateBook(command *CreateBookCommand) (*Book, error)
}

type service struct {
	repo BookRepository
}

func NewService(repo BookRepository) BookService {
	return &service{
		repo: repo,
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

	return &newBook, nil
}
