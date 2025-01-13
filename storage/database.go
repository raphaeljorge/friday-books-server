package storage

    import (
      "database/sql"
      "fmt"
      _ "github.com/lib/pq"
    )

    func InitDB() *sql.DB {
      connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
      )

      db, err := sql.Open("postgres", connStr)
      if err != nil {
        log.Fatal(err)
      }

      err = db.Ping()
      if err != nil {
        log.Fatal(err)
      }

      return db
    }

    func InitRedis() *redis.Client {
      return redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
      })
    }
