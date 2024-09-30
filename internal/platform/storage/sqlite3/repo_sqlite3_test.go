package sqlite3_test

import (
	"context"
	"database/sql"
	"library-api/internal"
	"library-api/internal/platform/storage/sqlite3"
	"library-api/kit/date"
	"os"
	"strings"
	"testing"
	"time"

	lorelai "github.com/UltiRequiem/lorelai/pkg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db                 *sql.DB
	transactionManager *sqlite3.TransactionManager
	repo               internal.BookRepository
}

const dbpath string = "testdb.sqlite3"
const schema string = `
CREATE TABLE IF NOT EXISTS books (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    publish_date DATE NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE  IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL
);

CREATE TABLE  IF NOT EXISTS book_categories (
    book_id TEXT,
    category_id INTEGER,
    FOREIGN KEY (book_id) REFERENCES books (id),
    FOREIGN KEY (category_id) REFERENCES categories (id),
    PRIMARY KEY (book_id, category_id)
);`

func TestRepo(t *testing.T) {
	if integrationEnabled := os.Getenv("RUN_INTEGRATION_TESTS"); integrationEnabled != "1" {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupSuite() {
	testDatabase := getTestDB()
	s.db = testDatabase
	s.transactionManager = sqlite3.NewTransactionManager(testDatabase)
	dateHandler := date.NewLocalHandler()
	s.repo = sqlite3.NewBookRepository(dateHandler)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	s.db.Close()
	os.Remove(dbpath)
}

func (s *RepositoryTestSuite) AfterTest(suiteName, testName string) {
	clearDatabase(s.db)
}

func (s *RepositoryTestSuite) TestFindAll() {
	_ = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		firstBooks, err := s.repo.FindAll(ctx)
		s.Assert().Nil(err)
		s.Assert().Empty(firstBooks)

		return nil
	})

	book1 := createTestBook()
	book2 := createTestBook()

	allBooks := []internal.Book{book1, book2}

	_ = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		s.repo.Create(ctx, book1)
		s.repo.Create(ctx, book2)
		booksStored, err := s.repo.FindAll(ctx)
		s.Assert().Nil(err)
		s.Assert().ElementsMatch(allBooks, booksStored)
		return nil
	})

}

func (s *RepositoryTestSuite) TestCreateAndFindById() {
	book1 := createTestBook()

	err := s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		return s.repo.Create(ctx, book1)
	})
	s.Require().NoError(err)

	_ = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		persistedBook, err := s.repo.FindByID(ctx, book1.ID)
		s.Assert().NoError(err)
		s.Assert().EqualValues(book1, persistedBook)
		return nil
	})
}

func (s *RepositoryTestSuite) TestFindByNonExistingId() {
	nonExistingId := uuid.New()

	_ = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		_, err := s.repo.FindByID(ctx, nonExistingId)
		s.Assert().ErrorIs(err, internal.ErrNotFound)
		return nil
	})
}

func (s *RepositoryTestSuite) TestCannotCreateDuplicatedBookId() {
	book1 := createTestBook()
	err := s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		return s.repo.Create(ctx, book1)
	})
	s.Require().NoError(err)

	err = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		return s.repo.Create(ctx, book1)
	})
	s.Require().Error(err)
}

func (s *RepositoryTestSuite) TestUpdateExistingBook() {
	book1 := createTestBook()

	err := s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		return s.repo.Create(ctx, book1)
	})
	s.Require().NoError(err)

	modifiedName, _ := internal.NewBookName("Test modified book")
	book1.Name = modifiedName

	modifiedDescription, _ := internal.NewBookDescription("Modificed description")
	book1.Description = modifiedDescription

	time := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	modifiedPublishedDate, _ := internal.NewBookPublishDate(time)
	book1.PublishDate = modifiedPublishedDate

	err = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		return s.repo.Update(ctx, book1)
	})
	s.Require().NoError(err)

	_ = s.transactionManager.RunInTransaction(context.Background(), func(ctx context.Context) error {
		persistedBook, err := s.repo.FindByID(ctx, book1.ID)
		s.Assert().NoError(err)
		s.Assert().EqualValues(book1, persistedBook)
		return nil
	})
}

func getTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	initSchema(db)

	return db
}

func createTestBook() internal.Book {
	bookId := uuid.New()

	book, _ := internal.BookFactory.NewBook(
		bookId,
		strings.TrimSpace(lorelai.LoremWords(2)),
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		strings.Split(strings.TrimSpace(lorelai.LoremWords(2)), " "),
		lorelai.Paragraph())

	return book
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

func initSchema(db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = transaction.Exec(schema)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}
