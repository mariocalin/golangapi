package book

import (
	"time"

	"github.com/google/uuid"
)

type BookId = uuid.UUID

type Book struct {
	ID          BookId    `json:"id"`
	Name        string    `json:"name"`
	PublishDate time.Time `json:"publish_date"`
	Categories  []string  `json:"categories"`
	Description string    `json:"description"`
}
