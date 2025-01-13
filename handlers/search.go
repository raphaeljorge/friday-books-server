package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "net/http"
      "strconv"
      "strings"
    )

    func SearchBooks(w http.ResponseWriter, r *http.Request) {
      query := r.URL.Query()
      page, _ := strconv.Atoi(query.Get("page"))
      if page < 1 {
        page = 1
      }

      limit, _ := strconv.Atoi(query.Get("limit"))
      if limit < 1 || limit > 100 {
        limit = 20
      }

      offset := (page - 1) * limit

      searchParams := map[string]string{
        "title":    query.Get("title"),
        "author":   query.Get("author"),
        "genre":    query.Get("genre"),
        "language": query.Get("language"),
      }

      var whereClauses []string
      var args []interface{}
      argCounter := 1

      for key, value := range searchParams {
        if value != "" {
          whereClauses = append(whereClauses, fmt.Sprintf("%s ILIKE $%d", key, argCounter))
          args = append(args, "%"+value+"%")
          argCounter++
        }
      }

      baseQuery := `
        SELECT id, title, authors, description, cover_image_url, average_rating 
        FROM books
      `

      if len(whereClauses) > 0 {
        baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
      }

      baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

      db := storage.GetDB()
      rows, err := db.Query(baseQuery, args...)
      if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }
      defer rows.Close()

      var books []models.Book
      for rows.Next() {
        var book models.Book
        err := rows.Scan(
          &book.ID,
          &book.Title,
          &book.Authors,
          &book.Description,
          &book.CoverImageURL,
          &book.AverageRating,
        )
        if err != nil {
          http.Error(w, "Error scanning results", http.StatusInternalServerError)
          return
        }
        books = append(books, book)
      }

      json.NewEncoder(w).Encode(map[string]interface{}{
        "page":       page,
        "limit":      limit,
        "total":      len(books),
        "results":    books,
      })
    }
