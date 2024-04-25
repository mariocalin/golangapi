package book

type BookRepository interface {
	FindAll() ([]Book, error)
	FindByID(id BookId) (*Book, error)
	Create(book *Book) error
	Update(book *Book) error
}
