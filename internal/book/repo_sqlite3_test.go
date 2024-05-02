package book

import (
	"database/sql"
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
	repo BookRepository
}

const dbpath string = "testdb.sqlite3"

func (suite *MySuite) SetupSuite() {
	testDatabase := getTestDB()
	suite.repo = NewSqlite3BookRepository(testDatabase)
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

	var allBooks []Book

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

	modifiedName, _ := NewName("Test modified book")
	book1.Name = modifiedName

	modifiedDescription, _ := NewDescription("Modificed description")
	book1.Description = modifiedDescription

	modifiedPublishedDate := NewPublishDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	book1.PublishDate = modifiedPublishedDate

	suite.repo.Update(book1)
	persistedBook, _ := suite.repo.FindByID(book1.ID)
	suite.Assert().EqualValues(book1, persistedBook)
}

func TestRepo(t *testing.T) {
	if integrationEnabled := os.Getenv("RUN_INTEGRATION_TESTS"); integrationEnabled != "1" {
		t.Skip("Skipping integration test")
	}

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

func createTestBook() *Book {
	bookId := uuid.New()

	return &Book{
		ID:          &bookId,
		Name:        &Name{strings.TrimSpace(lorelai.LoremWords(2))},
		PublishDate: &PublishDate{time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)},
		Categories:  &Categories{strings.Split(strings.TrimSpace(lorelai.LoremWords(2)), " ")},
		Description: &Description{lorelai.Paragraph()},
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
