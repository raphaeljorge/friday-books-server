CREATE TABLE users (
      id SERIAL PRIMARY KEY,
      username VARCHAR(255) UNIQUE NOT NULL,
      password VARCHAR(255) NOT NULL,
      email VARCHAR(255) UNIQUE NOT NULL,
      points INTEGER DEFAULT 0,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE books (
      id SERIAL PRIMARY KEY,
      title VARCHAR(255) NOT NULL,
      subtitle VARCHAR(255),
      authors JSONB,
      editors JSONB,
      translators JSONB,
      publisher VARCHAR(255),
      publication_date DATE,
      isbn_10 VARCHAR(10),
      isbn_13 VARCHAR(13),
      language VARCHAR(50),
      genres JSONB,
      tags JSONB,
      description TEXT,
      edition VARCHAR(50),
      page_count INTEGER,
      target_audience VARCHAR(50),
      series_name VARCHAR(255),
      volume_number INTEGER,
      awards JSONB,
      country_of_origin VARCHAR(100),
      cover_image_url VARCHAR(255),
      price NUMERIC(10, 2),
      currency VARCHAR(3),
      retailers JSONB,
      average_rating NUMERIC(3, 2),
      rating_count INTEGER,
      related_isbns JSONB,
      custom_notes TEXT,
      format VARCHAR(50),
      print_dimensions VARCHAR(50),
      weight NUMERIC(10, 2),
      illustrators JSONB,
      illustration_count INTEGER,
      dust_jacket VARCHAR(255),
      paper_type VARCHAR(50),
      ebook_format VARCHAR(50),
      ebook_size BIGINT,
      narrators JSONB,
      runtime INTEGER,
      audio_format VARCHAR(50),
      audio_size BIGINT,
      sample_audio_url VARCHAR(255),
      chapter_indexed BOOLEAN,
      bitrate INTEGER,
      download_url VARCHAR(255),
      audio_publisher VARCHAR(255),
      audio_rights_owner VARCHAR(255),
      audio_release_date DATE,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE reviews (
      id SERIAL PRIMARY KEY,
      book_id INTEGER REFERENCES books(id),
      user_id INTEGER REFERENCES users(id),
      rating INTEGER CHECK (rating BETWEEN 1 AND 5),
      comment TEXT,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE api_tokens (
      id SERIAL PRIMARY KEY,
      user_id INTEGER REFERENCES users(id),
      token VARCHAR(255) UNIQUE NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX idx_books_title ON books USING gin (to_tsvector('english', title));
    CREATE INDEX idx_books_authors ON books USING gin (authors);
    CREATE INDEX idx_books_genres ON books USING gin (genres);
