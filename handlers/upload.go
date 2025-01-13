package handlers

    import (
      "book-api/storage"
      "io"
      "net/http"
      "os"
      "path/filepath"
      "time"
    )

    func UploadCoverImage(w http.ResponseWriter, r *http.Request) {
      // Limit upload size to 5MB
      r.ParseMultipartForm(5 << 20)

      file, handler, err := r.FormFile("cover_image")
      if err != nil {
        http.Error(w, "Invalid file", http.StatusBadRequest)
        return
      }
      defer file.Close()

      // Validate file type
      allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/webp": true,
      }
      if !allowedTypes[handler.Header.Get("Content-Type")] {
        http.Error(w, "Invalid file type", http.StatusUnsupportedMediaType)
        return
      }

      // Create uploads directory if it doesn't exist
      uploadDir := "./uploads"
      if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        os.Mkdir(uploadDir, 0755)
      }

      // Generate unique filename
      ext := filepath.Ext(handler.Filename)
      newFilename := time.Now().Format("20060102150405") + ext
      filePath := filepath.Join(uploadDir, newFilename)

      // Save file
      dst, err := os.Create(filePath)
      if err != nil {
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
      }
      defer dst.Close()

      _, err = io.Copy(dst, file)
      if err != nil {
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
      }

      // Store file URL in database
      db := storage.GetDB()
      bookID := r.FormValue("book_id")
      _, err = db.Exec(
        "UPDATE books SET cover_image_url = $1 WHERE id = $2",
        "/uploads/"+newFilename,
        bookID,
      )
      if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
      }

      json.NewEncoder(w).Encode(map[string]string{
        "url": "/uploads/" + newFilename,
      })
    }
