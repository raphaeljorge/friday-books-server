package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "net/http"
      "time"
    )

    type AnalyticsResponse struct {
      TotalBooks        int `json:"total_books"`
      TotalUsers        int `json:"total_users"`
      TotalReviews      int `json:"total_reviews"`
      BooksAddedLast30  int `json:"books_added_last_30"`
      ActiveUsersLast30 int `json:"active_users_last_30"`
      PopularGenres     []GenreStats `json:"popular_genres"`
    }

    type GenreStats struct {
      Genre  string `json:"genre"`
      Count  int    `json:"count"`
    }

    func GetAnalytics(w http.ResponseWriter, r *http.Request) {
      db := storage.GetDB()
      var analytics AnalyticsResponse

      // Get total books
      db.QueryRow("SELECT COUNT(*) FROM books").Scan(&analytics.TotalBooks)

      // Get total users
      db.QueryRow("SELECT COUNT(*) FROM users").Scan(&analytics.TotalUsers)

      // Get total reviews
      db.QueryRow("SELECT COUNT(*) FROM reviews").Scan(&analytics.TotalReviews)

      // Get books added in last 30 days
      db.QueryRow(`
        SELECT COUNT(*) 
        FROM books 
        WHERE created_at >= $1`,
        time.Now().Add(-30*24*time.Hour),
      ).Scan(&analytics.BooksAddedLast30)

      // Get active users in last 30 days
      db.QueryRow(`
        SELECT COUNT(DISTINCT user_id) 
        FROM reviews 
        WHERE created_at >= $1`,
        time.Now().Add(-30*24*time.Hour),
      ).Scan(&analytics.ActiveUsersLast30)

      // Get popular genres
      rows, err := db.Query(`
        SELECT genre, COUNT(*) as count
        FROM books, jsonb_array_elements_text(genres) as genre
        GROUP BY genre
        ORDER BY count DESC
        LIMIT 5
      `)
      if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }
      defer rows.Close()

      for rows.Next() {
        var gs GenreStats
        rows.Scan(&gs.Genre, &gs.Count)
        analytics.PopularGenres = append(analytics.PopularGenres, gs)
      }

      json.NewEncoder(w).Encode(analytics)
    }
