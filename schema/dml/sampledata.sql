-- Insertar un nuevo libro en la tabla books
INSERT INTO books (id, name, publish_date, description) 
VALUES ('c9bf6d6d-6a5a-4d28-a6ba-35a5410e5fc3', 'El Quijote', '1605-01-01', 'El ingenioso hidalgo don Quijote de la Mancha, es una novela escrita por el español Miguel de Cervantes Saavedra.');

-- Insertar nuevas categorías si no existen ya en la tabla categories
INSERT INTO categories (category) VALUES ('Novela') ON CONFLICT DO NOTHING;
INSERT INTO categories (category) VALUES ('Literatura Española') ON CONFLICT DO NOTHING;

-- Insertar las relaciones entre el libro y las categorías en la tabla book_categories
INSERT INTO book_categories (book_id, category_id) 
SELECT 'c9bf6d6d-6a5a-4d28-a6ba-35a5410e5fc3', id FROM categories WHERE category IN ('Novela', 'Literatura Española');