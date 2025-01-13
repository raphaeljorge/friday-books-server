package handlers

    import (
      "database/sql"
      "net/http"
      "time"

      "github.com/go-redis/redis/v8"
    )

    type HealthResponse struct {
      Status    string `json:"status"`
      Database  string `json:"database"`
      Redis     string `json:"redis"`
      Timestamp string `json:"timestamp"`
    }

    func HealthCheck(w http.ResponseWriter, r *http.Request) {
      db := storage.GetDB()
      redisClient := storage.GetRedis()

      response := HealthResponse{
        Timestamp: time.Now().UTC().Format(time.RFC3339),
      }

      // Check database
      err := db.Ping()
      if err != nil {
        response.Status = "unhealthy"
        response.Database = "down"
      } else {
        response.Database = "up"
      }

      // Check Redis
      _, err = redisClient.Ping(ctx).Result()
      if err != nil {
        response.Status = "unhealthy"
        response.Redis = "down"
      } else {
        response.Redis = "up"
      }

      if response.Status == "" {
        response.Status = "healthy"
      }

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(response)
    }
