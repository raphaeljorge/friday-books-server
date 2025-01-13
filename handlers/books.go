package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "encoding/json"
      "net/http"
    )

    func CreateBook(w http.ResponseWriter, r *http.Request) {
      var book models.Book
      err := json.NewDecoder(r.Body).Decode(&book)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

      // Validate input
      if book.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
      }

      // Save to database
      db := storage.GetDB()
      // ... database insert logic ...

      w.WriteHeader(http.StatusCreated)
      json.NewEncoder(w).Encode(book)
    }

    // Implement other CRUD operations...
