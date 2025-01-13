package handlers

    import (
      "book-api/models"
      "book-api/storage"
      "context"
      "database/sql"
      "encoding/json"
      "log"
      "time"

      "github.com/hibiken/asynq"
    )

    type BackgroundTask struct {
      Type    string      `json:"type"`
      Payload interface{} `json:"payload"`
    }

    var redisOpt = asynq.RedisClientOpt{
      Addr:     os.Getenv("REDIS_ADDR"),
      Password: os.Getenv("REDIS_PASSWORD"),
    }

    func EnqueueTask(taskType string, payload interface{}) error {
      payloadBytes, err := json.Marshal(payload)
      if err != nil {
        return err
      }

      task := asynq.NewTask(taskType, payloadBytes)
      client := asynq.NewClient(redisOpt)
      defer client.Close()

      _, err = client.Enqueue(task)
      return err
    }

    func HandleBookIndexing(ctx context.Context, t *asynq.Task) error {
      var book models.Book
      if err := json.Unmarshal(t.Payload(), &book); err != nil {
        return err
      }

      // Index book in search engine
      db := storage.GetDB()
      _, err := db.Exec(`
        INSERT INTO book_search_index 
        (book_id, title, authors, description, genres)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (book_id) DO UPDATE SET
          title = EXCLUDED.title,
          authors = EXCLUDED.authors,
          description = EXCLUDED.description,
          genres = EXCLUDED.genres
      `, book.ID, book.Title, book.Authors, book.Description, book.Genres)

      return err
    }

    func StartBackgroundWorker() {
      srv := asynq.NewServer(
        redisOpt,
        asynq.Config{
          Concurrency: 10,
          Queues: map[string]int{
            "critical": 6,
            "default":  3,
            "low":      1,
          },
        },
      )

      mux := asynq.NewServeMux()
      mux.HandleFunc("book:index", HandleBookIndexing)

      if err := srv.Run(mux); err != nil {
        log.Fatalf("could not run server: %v", err)
      }
    }
