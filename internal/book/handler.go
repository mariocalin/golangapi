package book

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateBookRequest struct {
	Name        string   `json:"name" binding:"required"`
	PublishDate string   `json:"publish_date" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
	Description string   `json:"description" binding:"required"`
}

func getAllBooksHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := svc.GetBooks()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, createBookResources(books))
	}
}

func getBookByIdHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad id"})
			return
		}

		book, err := svc.GetBookByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if book == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
			return
		}

		c.JSON(http.StatusOK, createBookResource(book))
	}
}

func createBookHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req CreateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse publish date
		name, err := NewName(req.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
			return
		}

		description, err := NewDescription(req.Description)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid description"})
			return
		}

		categories, err := NewCategories(req.Categories)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid categories"})
			return
		}

		publishDate, err := time.Parse(time.RFC3339, req.PublishDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format"})
			return
		}

		// Create command
		command := CreateBookCommand{
			Name:        *name,
			Description: *description,
			Categories:  *categories,
			PublishDate: *NewPublishDate(publishDate)}

		// Add the book using the service
		book, err := svc.CreateBook(&command)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
			return
		}

		// Return the created book
		c.JSON(http.StatusCreated, createBookResource(book))
	}
}

func createBookResource(book *Book) BookResource {
	return BookResource{
		Id:          book.ID.String(),
		Name:        book.Name.Value(),
		PublishDate: book.PublishDate.Value().Format(time.DateOnly),
		Categories:  book.Categories.Value(),
		Description: book.Description.Value(),
	}
}

func createBookResources(books []Book) []BookResource {
	var resources []BookResource

	for _, book := range books {
		resources = append(resources, createBookResource(&book))
	}

	return resources
}

func RegisterHandlers(r *gin.Engine, svc BookService) {
	r.GET("/book", getAllBooksHandler(svc))
	r.GET("/book/:id", getBookByIdHandler(svc))
	r.POST("/book", createBookHandler(svc))
}

type BookResource struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	PublishDate string   `json:"publish_date"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
}
