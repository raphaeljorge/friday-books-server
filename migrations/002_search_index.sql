CREATE TABLE book_search_index (
      book_id UUID PRIMARY KEY REFERENCES books(id),
      title TEXT NOT NULL,
      authors TEXT[] NOT NULL,
      description TEXT,
      genres TEXT[] NOT NULL,
      tsvector TSVECTOR
    );

    CREATE INDEX idx_book_search ON book_search_index USING GIN (tsvector);

    CREATE OR REPLACE FUNCTION update_book_search_index() RETURNS TRIGGER AS $$
    BEGIN
      INSERT INTO book_search_index (book_id, title, authors, description, genres, tsvector)
      VALUES (NEW.id, NEW.title, NEW.authors, NEW.description, NEW.genres,
              setweight(to_tsvector('english', NEW.title), 'A') ||
              setweight(to_tsvector('english', array_to_string(NEW.authors, ' ')), 'B') ||
              setweight(to_tsvector('english', NEW.description), 'C'))
      ON CONFLICT (book_id) DO UPDATE SET
        title = EXCLUDED.title,
        authors = EXCLUDED.authors,
        description = EXCLUDED.description,
        genres = EXCLUDED.genres,
        tsvector = EXCLUDED.tsvector;
      RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER book_search_update
    AFTER INSERT OR UPDATE ON books
    FOR EACH ROW EXECUTE PROCEDURE update_book_search_index();
