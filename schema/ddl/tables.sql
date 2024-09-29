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