package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Id = uuid.UUID

type Book struct {
	ID          *Id
	Name        *Name
	PublishDate *PublishDate
	Categories  *Categories
	Description *Description
}

type Name struct {
	value string
}

func NewName(value string) (*Name, error) {
	if value == "" {
		return nil, errors.New("name cannot be empty")
	}

	if len(value) <= 0 || len(value) >= 80 {
		return nil, errors.New("name must be between 1 and 80 characters long")
	}

	return &Name{value: value}, nil
}

func (n *Name) Value() string {
	return n.value
}

type PublishDate struct {
	value time.Time
}

func NewPublishDate(value time.Time) *PublishDate {
	return &PublishDate{value: value}
}

func (p *PublishDate) Value() time.Time {
	return p.value
}

type Categories struct {
	value []string
}

func NewCategories(value []string) (*Categories, error) {
	if len(value) == 0 {
		return nil, errors.New("categories cannot be empty")
	}
	return &Categories{value: value}, nil
}

func (c *Categories) Value() []string {
	return c.value
}

type Description struct {
	value string
}

func NewDescription(value string) (*Description, error) {
	if value == "" {
		return nil, errors.New("description cannot be empty")
	}

	if len(value) <= 0 || len(value) >= 400 {
		return nil, errors.New("name must be between 1 and 80 characters long")
	}

	return &Description{value: value}, nil
}

func (d *Description) Value() string {
	return d.value
}
