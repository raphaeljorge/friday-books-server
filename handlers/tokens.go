package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "net/http"
      "time"

      "github.com/dgrijalva/jwt-go"
    )

    func GenerateAPIToken(w http.ResponseWriter, r *http.Request) {
      userID := r.Context().Value("user_id").(string)
      if userID == "" {
        http.Error(w, "Authentication required", http.StatusUnauthorized)
        return
      }

      claims := &Claims{
        UserID:   userID,
        StandardClaims: jwt.StandardClaims{
          ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
        },
      }

      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      tokenString, err := token.SignedString(jwtKey)
      if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
      }

      // Store token in database
      db := storage.GetDB()
      _, err = db.Exec(
        "INSERT INTO api_tokens (user_id, token) VALUES ($1, $2)",
        userID, tokenString,
      )

      if err != nil {
        http.Error(w, "Could not store token", http.StatusInternalServerError)
        return
      }

      json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
      })
    }

    func RevokeAPIToken(w http.ResponseWriter, r *http.Request) {
      token := mux.Vars(r)["token"]
      if token == "" {
        http.Error(w, "Token required", http.StatusBadRequest)
        return
      }

      db := storage.GetDB()
      _, err := db.Exec(
        "DELETE FROM api_tokens WHERE token = $1",
        token,
      )

      if err != nil {
        http.Error(w, "Could not revoke token", http.StatusInternalServerError)
        return
      }

      w.WriteHeader(http.StatusNoContent)
    }
