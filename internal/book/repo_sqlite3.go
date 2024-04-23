package book

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite3BookRepository struct {
	db *sql.DB
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


-- Insertar un nuevo libro en la tabla books
INSERT INTO books (id, name, publish_date, description) 
VALUES ('c9bf6d6d-6a5a-4d28-a6ba-35a5410e5fc3', 'El Quijote', '1605-01-01', 'El ingenioso hidalgo don Quijote de la Mancha, es una novela escrita por el español Miguel de Cervantes Saavedra.');

-- Insertar nuevas categorías si no existen ya en la tabla categories
INSERT INTO categories (category) VALUES ('Novela') ON CONFLICT DO NOTHING;
INSERT INTO categories (category) VALUES ('Literatura Española') ON CONFLICT DO NOTHING;

-- Insertar las relaciones entre el libro y las categorías en la tabla book_categories
INSERT INTO book_categories (book_id, category_id) 
SELECT 'c9bf6d6d-6a5a-4d28-a6ba-35a5410e5fc3', id FROM categories WHERE category IN ('Novela', 'Literatura Española');
`

func NewSqlite3BookRepository(dbpath string) BookRepository {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	if err := initSchema(db); err != nil {
		panic(err)
	}

	return &sqlite3BookRepository{db}
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
	SELECT b.id, b.name, b.publish_date, b.description, GROUP_CONCAT(c.category) 
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

		publishDate, _ := time.Parse(time.DateOnly, publishDateStr)

		bookRow.PublishDate = publishDate
		bookRow.Categories = strings.Split(categories, ",")

		book := toBook(&bookRow)
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *sqlite3BookRepository) FindByID(id BookId) (*Book, error) {
	var bookRow bookRow
	var publishDateStr string
	var categories string

	err := r.db.QueryRow(`
        SELECT b.id, b.name, b.publish_date, b.description, GROUP_CONCAT(c.category) 
        FROM books b 
            LEFT JOIN book_categories bc ON b.id = bc.book_id 
            LEFT JOIN categories c ON bc.category_id = c.id 
        WHERE b.id = ? 
        GROUP BY b.id`, id).Scan(&bookRow.Id, &bookRow.Name, &publishDateStr, &bookRow.Description, &categories)
	if err != nil {
		return nil, err
	}

	publishDate, err := time.Parse(time.DateOnly, publishDateStr)
	if err != nil {
		return nil, err
	}

	bookRow.PublishDate = publishDate
	bookRow.Categories = strings.Split(categories, ",")

	book := toBook(&bookRow)

	return &book, nil
}

func (r *sqlite3BookRepository) Create(book *Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}()

	_, err = tx.Exec("INSERT INTO books (id, name, publish_date, description) VALUES (?, ?, ?, ?)",
		book.ID.String(), book.Name.Value(), book.PublishDate.Value().Format(time.DateOnly), book.Description.Value())
	if err != nil {
		return err
	}

	for _, category := range book.Categories.value {
		var categoryID int
		err := tx.QueryRow("INSERT INTO categories (category) VALUES (?) RETURNING id", category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO book_categories (book_id, category_id) VALUES (?, ?)",
			book.ID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *sqlite3BookRepository) Update(book *Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}()

	_, err = tx.Exec("DELETE FROM book_categories WHERE book_id = ?", book.ID.String())
	if err != nil {
		return err
	}

	for _, category := range book.Categories.value {
		var categoryID int
		err := tx.QueryRow("INSERT INTO categories (category) VALUES (?) RETURNING id", category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO book_categories (book_id, category_id) VALUES (?, ?)",
			book.ID.String(), categoryID)
		if err != nil {
			return err
		}
	}

	return nil

}

func toBook(br *bookRow) Book {
	id, _ := uuid.Parse(br.Id)

	return Book{
		ID:          id,
		Name:        Name{value: br.Name},
		PublishDate: PublishDate{value: br.PublishDate},
		Description: Description{value: br.Description},
		Categories:  Categories{value: br.Categories},
	}
}
