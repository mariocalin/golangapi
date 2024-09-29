package bootstrap

import (
	"context"
	"fmt"
	"library-api/internal"
	"library-api/internal/book/create"
	"library-api/internal/book/findAll"
	"library-api/internal/book/findById"
	"library-api/internal/book/update"
	"library-api/internal/platform/event/memory"
	"library-api/internal/platform/server"
	"library-api/internal/platform/storage/sqlite3"
	"library-api/kit/date"
	"library-api/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

//	@title			Library Api
//	@version		1.0
//	@description	API for creating and retreiving books

//	@contact.name	Mario
//	@contact.url	http://example.org
//	@contact.email	mario.calin@mindcurv.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func Run() error {
	log.Info("Starting api")

	db := getDatabase()
	txManager := sqlite3.NewTransactionManager(db)
	dateHandler := date.NewLocalHandler()
	bookRepository := sqlite3.NewBookRepository(dateHandler)
	bookCreatedChannel := make(chan internal.BookCreated)
	bookEventPropagator := memory.NewChannelBookEventPropagator(bookCreatedChannel)

	bookEventConsumer := memory.NewChannelBookConsumer(bookCreatedChannel)

	createUseCase := create.NewUseCase(bookRepository, txManager, bookEventPropagator)
	findAllUseCase := findAll.NewUseCase(bookRepository, txManager)
	findByIdUseCase := findById.NewUseCase(bookRepository, txManager)
	updateUseCase := update.NewUseCase(bookRepository, txManager)

	bookAdapter := server.NewBookAdapter(createUseCase, findAllUseCase, findByIdUseCase, updateUseCase)

	bookEventConsumer.BindBookCreated(func(ctx context.Context, bc internal.BookCreated) {
		bok, _ := bookAdapter.GetBookByID(ctx, bc.Id.String())
		log.Info(fmt.Sprintf("A book has been created: %s", bok.String()))
	})

	go bookEventConsumer.StartConsuming(context.Background())
	defer bookEventConsumer.Close(context.Background())
	defer bookEventPropagator.Close(context.Background())

	router := gin.Default()
	bookController := server.NewBookController(bookAdapter, dateHandler)
	registerHandlers(router, bookController)

	// ---- STATUS ----
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go startServer(router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Aplicación stopped")

	return nil
}

func startServer(gin *gin.Engine) {
	// Iniciar servidor en rutina
	if err := gin.Run(":8080"); err != nil {
		panic(fmt.Sprintf("Error running HTTP Server: %s", err.Error()))
	}
}

func registerHandlers(r *gin.Engine, bookController server.BookController) {
	r.GET("/book", bookController.GetAllBooks)
	r.GET("/book/:id", bookController.GetBookById)
	r.POST("/book", bookController.CreateBook)
}
