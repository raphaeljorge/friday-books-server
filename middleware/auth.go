package middleware

    import (
      "net/http"
      "strings"
      "github.com/dgrijalva/jwt-go"
    )

    func JWTAuth(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
          http.Error(w, "Authorization header required", http.StatusUnauthorized)
          return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
          }
          return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
          http.Error(w, "Invalid token", http.StatusUnauthorized)
          return
        }

        next.ServeHTTP(w, r)
      })
    }
