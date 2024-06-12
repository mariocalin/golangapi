//go:build !unit || integration

package sqlite

import (
	"database/sql"
	"library-api/common"
	"library-api/internal/book/application"
	"library-api/internal/book/domain"
	"os"
	"strings"
	"testing"
	"time"

	lorelai "github.com/UltiRequiem/lorelai/pkg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type MySuite struct {
	suite.Suite
	db   *sql.DB
	repo application.BookRepository
}

const dbpath string = "testdb.sqlite3"

func (suite *MySuite) SetupSuite() {
	testDatabase := getTestDB()
	dateHandler := common.NewDateHandler()
	suite.repo = NewSqlite3BookRepository(testDatabase, dateHandler)
	suite.db = testDatabase
}

func (suite *MySuite) TearDownSuite() {
	suite.db.Close()
	os.Remove(dbpath)
}

func (suite *MySuite) AfterTest(suiteName, testName string) {
	clearDatabase(suite.db)
}

func (suite *MySuite) TestFindAll() {
	firstBooks, err := suite.repo.FindAll()
	suite.Assert().Nil(err)
	suite.Assert().Empty(firstBooks)

	var allBooks []domain.Book

	book1 := createTestBook()
	suite.repo.Create(book1)
	allBooks = append(allBooks, *book1)

	book2 := createTestBook()
	suite.repo.Create(book2)
	allBooks = append(allBooks, *book2)

	books, _ := suite.repo.FindAll()
	suite.Assert().ElementsMatch(books, allBooks)
}

func (suite *MySuite) TestCreateAndFindById() {
	book1 := createTestBook()
	suite.repo.Create(book1)
	persistedBook, err := suite.repo.FindByID(book1.ID)
	suite.Assert().Nil(err)
	suite.Assert().EqualValues(book1, persistedBook)
}

func (suite *MySuite) TestFindByNonExistingId() {
	nonExistingId := uuid.New()
	book, err := suite.repo.FindByID(&nonExistingId)
	suite.Assert().Nil(book)
	suite.Assert().NotNil(err)
}

func (suite *MySuite) TestCannotCreateDuplicatedBookId() {
	book1 := createTestBook()
	suite.repo.Create(book1)
	err := suite.repo.Create(book1)
	suite.Assert().NotNil(err)
}

func (suite *MySuite) TestUpdateExistingBook() {
	book1 := createTestBook()
	suite.repo.Create(book1)

	modifiedName, _ := domain.NewName("Test modified book")
	book1.Name = modifiedName

	modifiedDescription, _ := domain.NewDescription("Modificed description")
	book1.Description = modifiedDescription

	modifiedPublishedDate := domain.NewPublishDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	book1.PublishDate = modifiedPublishedDate

	suite.repo.Update(book1)
	persistedBook, _ := suite.repo.FindByID(book1.ID)
	suite.Assert().EqualValues(book1, persistedBook)
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(MySuite))
}

func getTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	initSchema(db)

	return db
}

func createTestBook() *domain.Book {
	bookId := uuid.New()

	name, _ := domain.NewName(strings.TrimSpace(lorelai.LoremWords(2)))
	publishDate := domain.NewPublishDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	categories, _ := domain.NewCategories(strings.Split(strings.TrimSpace(lorelai.LoremWords(2)), " "))
	description, _ := domain.NewDescription(lorelai.Paragraph())

	return &domain.Book{
		ID:          &bookId,
		Name:        name,
		PublishDate: publishDate,
		Categories:  categories,
		Description: description,
	}
}

func clearDatabase(db *sql.DB) {
	transaction, _ := db.Begin()

	_, err := db.Exec("DELETE FROM book_categories;")
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("DELETE FROM categories;")
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("DELETE FROM books;")
	if err != nil {
		panic(err.Error())
	}

	transaction.Commit()
}
