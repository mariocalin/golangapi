package book

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BookId = uuid.UUID

type Name struct {
	value string
}

// NewName crea un nuevo objeto Name.
func NewName(value string) (*Name, error) {
	if value == "" {
		return nil, errors.New("name cannot be empty")
	}

	if len(value) <= 0 || len(value) >= 80 {
		return nil, errors.New("name must be between 1 and 80 characters long")
	}

	return &Name{value: value}, nil
}

// Value devuelve el valor del nombre.
func (n *Name) Value() string {
	return n.value
}

// PublishDate representa la fecha de publicación de un libro.
type PublishDate struct {
	value time.Time
}

// NewPublishDate crea un nuevo objeto PublishDate.
func NewPublishDate(value time.Time) *PublishDate {
	return &PublishDate{value: value}
}

// Value devuelve el valor de la fecha de publicación.
func (p *PublishDate) Value() time.Time {
	return p.value
}

// Categories representa las categorías de un libro.
type Categories struct {
	value []string
}

// NewCategories crea un nuevo objeto Categories.
func NewCategories(value []string) (*Categories, error) {
	if len(value) == 0 {
		return nil, errors.New("categories cannot be empty")
	}
	return &Categories{value: value}, nil
}

// Value devuelve el valor de las categorías.
func (c *Categories) Value() []string {
	return c.value
}

// Description representa la descripción de un libro.
type Description struct {
	value string
}

// NewDescription crea un nuevo objeto Description.
func NewDescription(value string) (*Description, error) {
	if value == "" {
		return nil, errors.New("description cannot be empty")
	}

	if len(value) <= 0 || len(value) >= 400 {
		return nil, errors.New("name must be between 1 and 80 characters long")
	}

	return &Description{value: value}, nil
}

// Value devuelve el valor de la descripción.
func (d *Description) Value() string {
	return d.value
}

type Book struct {
	ID          BookId
	Name        Name
	PublishDate PublishDate
	Categories  Categories
	Description Description
}
