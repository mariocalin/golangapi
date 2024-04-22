package book

type BookService interface {
	GetBooks() ([]Book, error)
	GetBookByID(id BookId) (*Book, error)
	CreateBook(book *Book) error
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

func (s *service) CreateBook(book *Book) error {
	return s.repo.Create(book)
}
