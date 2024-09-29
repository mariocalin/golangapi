package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BookId = uuid.UUID

func NewBookId() BookId {
	return uuid.New()
}

type Book struct {
	ID          BookId
	Name        BookName
	PublishDate BookPublishDate
	Categories  BookCategories
	Description BookDescription
}

type BookName struct {
	value string
}

var ErrNameEmpty = errors.New("name cannot be empty")
var ErrNameLength = errors.New("name must be between 1 and 80 characters long")

func NewBookName(value string) (BookName, error) {
	if value == "" {
		return BookName{}, ErrNameEmpty
	}

	if len(value) <= 0 || len(value) >= 80 {
		return BookName{}, ErrNameLength
	}

	return BookName{value: value}, nil
}

func (n BookName) Value() string {
	return n.value
}

type BookPublishDate struct {
	value time.Time
}

var ErrPublishDateEmpty = errors.New("publish date cannot be empty")

func NewBookPublishDate(value time.Time) (BookPublishDate, error) {
	if value.IsZero() {
		return BookPublishDate{}, ErrPublishDateEmpty
	}

	return BookPublishDate{value: value}, nil
}

func (p BookPublishDate) Value() time.Time {
	return p.value
}

type BookCategories struct {
	value []string
}

var ErrCategoriesEmpty = errors.New("categories cannot be empty")

func NewBookCategories(value []string) (BookCategories, error) {
	if value == nil {
		return BookCategories{}, ErrCategoriesEmpty
	}

	return BookCategories{value: value}, nil
}

func (c BookCategories) Value() []string {
	return c.value
}

type BookDescription struct {
	value string
}

var ErrDescriptionEmpty = errors.New("description cannot be empty")
var ErrDescriptionLength = errors.New("description must be between 1 and 400 characters long")

func NewBookDescription(value string) (BookDescription, error) {
	if value == "" {
		return BookDescription{}, ErrDescriptionEmpty
	}

	if len(value) <= 0 || len(value) >= 400 {
		return BookDescription{}, ErrDescriptionLength
	}

	return BookDescription{value: value}, nil
}

func (d BookDescription) Value() string {
	return d.value
}

func (b Book) String() string {
	return fmt.Sprintf("Book[ID=%s, Name=%s, PublishDate=%s, Categories=%v, Description=%s]",
		b.ID, b.Name.Value(), b.PublishDate.Value().Format(time.RFC3339), b.Categories.Value(), b.Description.Value())
}
