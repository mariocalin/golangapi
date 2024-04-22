package book

import (
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func getBookByIdHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		book, err := svc.GetBookByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if book == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusOK, book)
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
		command := CreateBookCommand{Name: *name, Description: *description, Categories: *categories, PublishDate: *NewPublishDate(publishDate)}

		// Add the book using the service
		book, err := svc.CreateBook(&command)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
			return
		}

		// Return the created book
		c.JSON(http.StatusCreated, gin.H{
			"id":           book.ID,
			"name":         book.Name.Value(),
			"publish_date": book.PublishDate.Value(),
			"categories":   book.Categories.Value(),
			"description":  book.Description.Value(),
		})
	}
}

func RegisterHandlers(r *gin.Engine, svc BookService) {
	r.GET("/book", getAllBooksHandler(svc))
	r.GET("/book/:id", getBookByIdHandler(svc))
	r.POST("/book", createBookHandler(svc))
}
