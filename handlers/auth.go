package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "database/sql"
      "encoding/json"
      "errors"
      "net/http"
      "time"

      "github.com/dgrijalva/jwt-go"
      "golang.org/x/crypto/bcrypt"
    )

    var jwtKey = []byte(os.Getenv("JWT_SECRET"))

    type Credentials struct {
      Username string `json:"username"`
      Password string `json:"password"`
    }

    type Claims struct {
      Username string `json:"username"`
      UserID   string `json:"user_id"`
      jwt.StandardClaims
    }

    func Register(w http.ResponseWriter, r *http.Request) {
      var user models.User
      err := json.NewDecoder(r.Body).Decode(&user)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
      if err != nil {
        http.Error(w, "Could not hash password", http.StatusInternalServerError)
        return
      }

      user.Password = string(hashedPassword)
      db := storage.GetDB()

      // Insert user into database
      err = db.QueryRow(
        "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id",
        user.Username, user.Password, user.Email,
      ).Scan(&user.ID)

      if err != nil {
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
      }

      w.WriteHeader(http.StatusCreated)
      json.NewEncoder(w).Encode(user)
    }

    func Login(w http.ResponseWriter, r *http.Request) {
      var creds Credentials
      err := json.NewDecoder(r.Body).Decode(&creds)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

      db := storage.GetDB()
      var user models.User
      err = db.QueryRow(
        "SELECT id, username, password FROM users WHERE username = $1",
        creds.Username,
      ).Scan(&user.ID, &user.Username, &user.Password)

      if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
          http.Error(w, "Invalid credentials", http.StatusUnauthorized)
          return
        }
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }

      err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
      if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
      }

      expirationTime := time.Now().Add(24 * time.Hour)
      claims := &Claims{
        Username: user.Username,
        UserID:   user.ID,
        StandardClaims: jwt.StandardClaims{
          ExpiresAt: expirationTime.Unix(),
        },
      }

      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      tokenString, err := token.SignedString(jwtKey)
      if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
      }

      http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
      })

      json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
      })
    }
