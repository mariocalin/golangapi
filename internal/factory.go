package internal

import "time"

var BookFactory = bookFactory{}

type bookFactory struct{}

func (bookFactory) NewBook(id BookId, name string, publishDate time.Time, categories []string, description string) (Book, error) {
	bookName, err := NewBookName(name)
	if err != nil {
		return Book{}, err
	}

	bookPublishDate, err := NewBookPublishDate(publishDate)
	if err != nil {
		return Book{}, err
	}

	bookCategories, err := NewBookCategories(categories)
	if err != nil {
		return Book{}, err
	}

	bookDescription, err := NewBookDescription(description)
	if err != nil {
		return Book{}, err
	}

	return Book{
		ID:          id,
		Name:        bookName,
		PublishDate: bookPublishDate,
		Categories:  bookCategories,
		Description: bookDescription,
	}, nil
}

func (f bookFactory) CreateBook(name string, publishDate time.Time, categories []string, description string) (Book, error) {
	return f.NewBook(NewBookId(), name, publishDate, categories, description)
}
