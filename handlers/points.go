package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "net/http"
      "strconv"
    )

    func GetUserPoints(w http.ResponseWriter, r *http.Request) {
      userID := mux.Vars(r)["id"]
      if userID == "" {
        http.Error(w, "User ID required", http.StatusBadRequest)
        return
      }

      db := storage.GetDB()
      var points int
      err := db.QueryRow(
        "SELECT points FROM users WHERE id = $1",
        userID,
      ).Scan(&points)

      if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
          http.Error(w, "User not found", http.StatusNotFound)
          return
        }
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }

      json.NewEncoder(w).Encode(map[string]int{
        "points": points,
      })
    }

    func AwardPoints(userID string, points int) error {
      db := storage.GetDB()
      _, err := db.Exec(
        "UPDATE users SET points = points + $1 WHERE id = $2",
        points, userID,
      )
      return err
    }
