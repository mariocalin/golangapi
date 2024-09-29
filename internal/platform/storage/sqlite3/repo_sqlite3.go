package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"library-api/internal"
	"library-api/kit/date"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite3BookRepository struct {
	dateHandler date.Handler
}

type bookRow struct {
	Id          string
	Name        string
	PublishDate time.Time
	Description string
	Categories  []string
}

func NewBookRepository(dateHandler date.Handler) internal.BookRepository {
	return &sqlite3BookRepository{dateHandler: dateHandler}
}

func extractTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(internal.TransactionKey).(*sql.Tx)
	if !ok {
		return nil, internal.ErrTransactionNotFound
	}

	return tx, nil
}

func (r sqlite3BookRepository) FindAll(ctx context.Context) ([]internal.Book, error) {
	tx, err := extractTransaction(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `
	SELECT b.id, b.name, DATE(b.publish_date), b.description, GROUP_CONCAT(c.category) 
	FROM books b 
		LEFT JOIN book_categories bc ON b.id = bc.book_id 
		LEFT JOIN categories c ON bc.category_id = c.id 
	GROUP BY b.id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []internal.Book

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

		book, err := toBook(bookRow)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r sqlite3BookRepository) FindByID(ctx context.Context, id internal.BookId) (internal.Book, error) {
	tx, err := extractTransaction(ctx)
	if err != nil {
		return internal.Book{}, err
	}

	stmt, err := tx.Prepare(`
        SELECT b.id, b.name, DATE(b.publish_date), b.description, GROUP_CONCAT(c.category) 
        FROM books b 
            LEFT JOIN book_categories bc ON b.id = bc.book_id 
            LEFT JOIN categories c ON bc.category_id = c.id 
        WHERE b.id = ? 
        GROUP BY b.id`)
	if err != nil {
		return internal.Book{}, err
	}

	defer stmt.Close()

	var bookRow bookRow
	var publishDateStr string
	var categories string

	// Ejecutar la consulta preparada
	err = stmt.QueryRow(id.String()).Scan(&bookRow.Id, &bookRow.Name, &publishDateStr, &bookRow.Description, &categories)
	if err != nil {
		return internal.Book{}, err
	}

	publishDate, err := time.ParseInLocation(time.DateOnly, publishDateStr, time.Local)
	if err != nil {
		return internal.Book{}, err
	}

	bookRow.PublishDate = publishDate
	bookRow.Categories = strings.Split(categories, ",")

	return toBook(bookRow)
}

func (r sqlite3BookRepository) Create(ctx context.Context, book internal.Book) error {
	transaction, err := extractTransaction(ctx)
	if err != nil {
		return err
	}

	// Preparar la inserción del libro
	bookStmt, err := transaction.Prepare("INSERT INTO books (id, name, publish_date, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(book.ID.String(), book.Name.Value(), book.PublishDate.Value().Format(time.DateOnly),
		book.Description.Value())
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

	for _, category := range book.Categories.Value() {
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

func (r sqlite3BookRepository) Update(ctx context.Context, book internal.Book) error {
	transaction, err := extractTransaction(ctx)
	if err != nil {
		return err
	}

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
	for _, category := range book.Categories.Value() {
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

func toBook(br bookRow) (internal.Book, error) {
	id, err := uuid.Parse(br.Id)
	if err != nil {
		return internal.Book{}, fmt.Errorf("could not parse book ID: %w", err)
	}

	name, err := internal.NewBookName(br.Name)
	if err != nil {
		return internal.Book{}, fmt.Errorf("could not create book name: %w", err)
	}

	publishDate, err := internal.NewBookPublishDate(br.PublishDate)
	if err != nil {
		return internal.Book{}, fmt.Errorf("could not create book publish date: %w", err)
	}

	description, err := internal.NewBookDescription(br.Description)
	if err != nil {
		return internal.Book{}, fmt.Errorf("could not create book description: %w", err)
	}

	categories, err := internal.NewBookCategories(br.Categories)
	if err != nil {
		return internal.Book{}, fmt.Errorf("could not create book categories: %w", err)
	}

	return internal.Book{
		ID:          id,
		Name:        name,
		PublishDate: publishDate,
		Description: description,
		Categories:  categories,
	}, nil
}
