package handlers

    import (
      "book-api/storage"
      "context"
      "encoding/json"
      "net/http"
      "time"

      "github.com/go-redis/redis/v8"
    )

    var ctx = context.Background()

    func GetCachedBooks(w http.ResponseWriter, r *http.Request) {
      cacheKey := "books:all"
      redisClient := storage.GetRedis()

      // Try to get from cache
      cachedData, err := redisClient.Get(ctx, cacheKey).Result()
      if err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedData))
        return
      }

      // If not in cache, get from database
      db := storage.GetDB()
      rows, err := db.Query(`
        SELECT id, title, authors, description, cover_image_url, average_rating
        FROM books
        ORDER BY created_at DESC
        LIMIT 100
      `)
      if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }
      defer rows.Close()

      var books []models.Book
      for rows.Next() {
        var book models.Book
        rows.Scan(
          &book.ID,
          &book.Title,
          &book.Authors,
          &book.Description,
          &book.CoverImageURL,
          &book.AverageRating,
        )
        books = append(books, book)
      }

      // Cache the result
      jsonData, _ := json.Marshal(books)
      redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

      w.Header().Set("Content-Type", "application/json")
      w.Write(jsonData)
    }
