package storage

    import (
      "database/sql"
      "os"
      "testing"

      _ "github.com/lib/pq"
    )

    func InitTestDB() *sql.DB {
      connStr := "host=localhost port=5432 user=postgres password=postgres dbname=books_test sslmode=disable"
      db, err := sql.Open("postgres", connStr)
      if err != nil {
        panic(err)
      }

      err = db.Ping()
      if err != nil {
        panic(err)
      }

      // Run migrations
      migration, err := os.ReadFile("../migrations/001_initial_schema.sql")
      if err != nil {
        panic(err)
      }

      _, err = db.Exec(string(migration))
      if err != nil {
        panic(err)
      }

      return db
    }
