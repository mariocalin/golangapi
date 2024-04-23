package main

import (
	"library-api/internal/book"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()

	// Book context
	bookRepo := book.NewSqlite3BookRepository(uuid.New().String() + ".sqlite3")
	bookSvc := book.NewService(bookRepo)
	book.RegisterHandlers(router, bookSvc)

	// Ping context
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Iniciar servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
