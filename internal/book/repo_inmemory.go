package book

import "errors"

type memoryBookRepository struct {
	books map[BookId]Book
}

func NewInMemoryBookRepository() BookRepository {
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
		return nil, nil
	}
	return &book, nil
}

func (r *memoryBookRepository) Create(book *Book) error {
	_, exists := r.books[book.ID]
	if exists {
		return errors.New("book already exists")
	}

	r.books[book.ID] = *book
	return nil
}

func (r *memoryBookRepository) Update(book *Book) error {
	r.books[book.ID] = *book
	return nil
}
