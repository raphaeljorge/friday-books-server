package tests

    import (
      "book-api/handlers"
      "book-api/models"
      "book-api/storage"
      "bytes"
      "encoding/json"
      "net/http"
      "net/http/httptest"
      "testing"

      "github.com/stretchr/testify/assert"
    )

    func TestCreateBook(t *testing.T) {
      db := storage.InitTestDB()
      defer db.Close()

      book := models.Book{
        Title:       "Test Book",
        Authors:     []string{"Author 1"},
        Description: "Test Description",
      }

      jsonData, _ := json.Marshal(book)
      req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonData))
      req.Header.Set("Content-Type", "application/json")

      rr := httptest.NewRecorder()
      handler := http.HandlerFunc(handlers.CreateBook)

      handler.ServeHTTP(rr, req)

      assert.Equal(t, http.StatusCreated, rr.Code)

      var responseBook models.Book
      json.Unmarshal(rr.Body.Bytes(), &responseBook)

      assert.Equal(t, book.Title, responseBook.Title)
      assert.Equal(t, book.Authors, responseBook.Authors)
      assert.NotEmpty(t, responseBook.ID)
    }

    // Add more tests...
