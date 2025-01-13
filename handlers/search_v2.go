package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "net/http"
      "strconv"
      "strings"

      "github.com/lib/pq"
    )

    type SearchResponse struct {
      Results    []models.Book `json:"results"`
      Total      int           `json:"total"`
      Page       int           `json:"page"`
      TotalPages int           `json:"total_pages"`
      Facets     SearchFacets  `json:"facets"`
    }

    type SearchFacets struct {
      Genres    []FacetCount `json:"genres"`
      Languages []FacetCount `json:"languages"`
      Formats   []FacetCount `json:"formats"`
    }

    type FacetCount struct {
      Name  string `json:"name"`
      Count int    `json:"count"`
    }

    func SearchBooksV2(w http.ResponseWriter, r *http.Request) {
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
        "query":    query.Get("q"),
        "genre":    query.Get("genre"),
        "language": query.Get("language"),
        "format":   query.Get("format"),
      }

      var whereClauses []string
      var args []interface{}
      argCounter := 1

      if searchParams["query"] != "" {
        whereClauses = append(whereClauses, `
          to_tsvector('english', title) || 
          to_tsvector('english', description) @@ plainto_tsquery($1)
        `)
        args = append(args, searchParams["query"])
        argCounter++
      }

      if searchParams["genre"] != "" {
        whereClauses = append(whereClauses, "$"+strconv.Itoa(argCounter)+" = ANY(genres)")
        args = append(args, searchParams["genre"])
        argCounter++
      }

      if searchParams["language"] != "" {
        whereClauses = append(whereClauses, "language = $"+strconv.Itoa(argCounter))
        args = append(args, searchParams["language"])
        argCounter++
      }

      if searchParams["format"] != "" {
        whereClauses = append(whereClauses, "format = $"+strconv.Itoa(argCounter))
        args = append(args, searchParams["format"])
        argCounter++
      }

      baseQuery := `
        SELECT id, title, authors, description, cover_image_url, average_rating,
               language, format, genres
        FROM books
      `

      countQuery := "SELECT COUNT(*) FROM books"
      if len(whereClauses) > 0 {
        whereClause := " WHERE " + strings.Join(whereClauses, " AND ")
        baseQuery += whereClause
        countQuery += whereClause
      }

      baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

      db := storage.GetDB()

      // Get total count
      var total int
      err := db.QueryRow(countQuery, args...).Scan(&total)
      if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }

      // Get results
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
          &book.Language,
          &book.Format,
          pq.Array(&book.Genres),
        )
        if err != nil {
          http.Error(w, "Error scanning results", http.StatusInternalServerError)
          return
        }
        books = append(books, book)
      }

      // Get facets
      facets := getSearchFacets(db)

      response := SearchResponse{
        Results:    books,
        Total:      total,
        Page:       page,
        TotalPages: (total + limit - 1) / limit,
        Facets:     facets,
      }

      json.NewEncoder(w).Encode(response)
    }

    func getSearchFacets(db *sql.DB) SearchFacets {
      var facets SearchFacets

      // Get genre facets
      rows, _ := db.Query(`
        SELECT genre, COUNT(*) as count
        FROM books, jsonb_array_elements_text(genres) as genre
        GROUP BY genre
        ORDER BY count DESC
        LIMIT 10
      `)
      defer rows.Close()

      for rows.Next() {
        var fc FacetCount
        rows.Scan(&fc.Name, &fc.Count)
        facets.Genres = append(facets.Genres, fc)
      }

      // Get language facets
      rows, _ = db.Query(`
        SELECT language, COUNT(*) as count
        FROM books
        GROUP BY language
        ORDER BY count DESC
        LIMIT 5
      `)
      defer rows.Close()

      for rows.Next() {
        var fc FacetCount
        rows.Scan(&fc.Name, &fc.Count)
        facets.Languages = append(facets.Languages, fc)
      }

      // Get format facets
      rows, _ = db.Query(`
        SELECT format, COUNT(*) as count
        FROM books
        WHERE format IS NOT NULL
        GROUP BY format
        ORDER BY count DESC
        LIMIT 5
      `)
      defer rows.Close()

      for rows.Next() {
        var fc FacetCount
        rows.Scan(&fc.Name, &fc.Count)
        facets.Formats = append(facets.Formats, fc)
      }

      return facets
    }
