package book

import (
	"database/sql"
	"library-api/common"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite3BookRepository struct {
	db          *sql.DB
	dateHandler *common.DateHandler
}

type bookRow struct {
	Id          string
	Name        string
	PublishDate time.Time
	Description string
	Categories  []string
}

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
);
`

func NewSqlite3BookRepository(db *sql.DB, dateHandler *common.DateHandler) BookRepository {
	if err := initSchema(db); err != nil {
		panic(err)
	}

	return &sqlite3BookRepository{db: db, dateHandler: dateHandler}
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

func (r *sqlite3BookRepository) FindAll() ([]Book, error) {
	rows, err := r.db.Query(`
	SELECT b.id, b.name, DATE(b.publish_date), b.description, GROUP_CONCAT(c.category) 
	FROM books b 
		LEFT JOIN book_categories bc ON b.id = bc.book_id 
		LEFT JOIN categories c ON bc.category_id = c.id 
	GROUP BY b.id`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var bookRow bookRow
		var publishDateStr string
		var categories string

		if err := rows.Scan(&bookRow.Id, &bookRow.Name, &publishDateStr, &bookRow.Description, &categories); err != nil {
			return nil, err
		}

		publishDate, _ := time.ParseInLocation(time.DateOnly, publishDateStr, time.Local)

		bookRow.PublishDate = publishDate.Local()
		bookRow.Categories = strings.Split(categories, ",")

		book := toBook(&bookRow)
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *sqlite3BookRepository) FindByID(id *BookId) (*Book, error) {
	var bookRow bookRow
	var publishDateStr string
	var categories string

	stmt, err := r.db.Prepare(`
        SELECT b.id, b.name, DATE(b.publish_date), b.description, GROUP_CONCAT(c.category) 
        FROM books b 
            LEFT JOIN book_categories bc ON b.id = bc.book_id 
            LEFT JOIN categories c ON bc.category_id = c.id 
        WHERE b.id = ? 
        GROUP BY b.id`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	// Ejecutar la consulta preparada
	err = stmt.QueryRow(id.String()).Scan(&bookRow.Id, &bookRow.Name, &publishDateStr, &bookRow.Description, &categories)
	if err != nil {
		return nil, err
	}

	publishDate, err := time.ParseInLocation(time.DateOnly, publishDateStr, time.Local)
	if err != nil {
		return nil, err
	}

	bookRow.PublishDate = publishDate
	bookRow.Categories = strings.Split(categories, ",")

	book := toBook(&bookRow)

	return &book, nil
}

func (r *sqlite3BookRepository) Create(book *Book) error {
	transaction, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() error {
		if err != nil {
			transaction.Rollback()
			return err
		}
		transaction.Commit()
		return nil
	}()

	// Preparar la inserción del libro
	bookStmt, err := transaction.Prepare("INSERT INTO books (id, name, publish_date, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(book.ID.String(), book.Name.Value(), book.PublishDate.Value().Format(time.DateOnly), book.Description.Value())
	if err != nil {
		return err
	}

	// Preparar la inserción de categorías
	categoryStmt, err := transaction.Prepare("INSERT INTO categories (category) VALUES (?) RETURNING id")
	if err != nil {
		return err
	}
	defer categoryStmt.Close()

	linkStmt, err := transaction.Prepare("INSERT INTO book_categories (book_id, category_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer linkStmt.Close()

	for _, category := range book.Categories.value {
		var categoryID int
		err := categoryStmt.QueryRow(category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = linkStmt.Exec(book.ID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *sqlite3BookRepository) Update(book *Book) error {
	transaction, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
			return
		}
		transaction.Commit()
	}()

	bookStmt, err := transaction.Prepare("UPDATE books SET name = ?, publish_date = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(book.Name.Value(), book.PublishDate.Value().Format(time.DateOnly), book.Description.Value(), book.ID.String())
	if err != nil {
		return err
	}

	// Eliminar todas las categorías asociadas al libro
	_, err = transaction.Exec("DELETE FROM book_categories WHERE book_id = ?", book.ID.String())
	if err != nil {
		return err
	}

	// Preparar la inserción de categorías
	categoryStmt, err := transaction.Prepare("INSERT INTO categories (category) VALUES (?) RETURNING id")
	if err != nil {
		return err
	}
	defer categoryStmt.Close()

	linkStmt, err := transaction.Prepare("INSERT INTO book_categories (book_id, category_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer linkStmt.Close()

	// Insertar las nuevas categorías
	for _, category := range book.Categories.value {
		var categoryID int
		err := categoryStmt.QueryRow(category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = linkStmt.Exec(book.ID.String(), categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func toBook(br *bookRow) Book {
	id, _ := uuid.Parse(br.Id)

	return Book{
		ID:          &id,
		Name:        &Name{value: br.Name},
		PublishDate: &PublishDate{value: br.PublishDate},
		Description: &Description{value: br.Description},
		Categories:  &Categories{value: br.Categories},
	}
}
