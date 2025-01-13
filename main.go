package main

    import (
      "book-api/handlers"
      "book-api/middleware"
      "book-api/storage"
      "log"
      "net/http"
      "os"

      "github.com/gorilla/mux"
      "github.com/joho/godotenv"
    )

    func main() {
      err := godotenv.Load()
      if err != nil {
        log.Fatal("Error loading .env file")
      }

      db := storage.InitDB()
      redisClient := storage.InitRedis()
      defer db.Close()

      // Start background worker
      go handlers.StartBackgroundWorker()

      r := mux.NewRouter()

      // Middleware
      r.Use(middleware.Logging)
      r.Use(middleware.Recovery)
      r.Use(middleware.SecurityHeaders)
      r.Use(middleware.CORS)

      // Static files
      r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

      // Health check
      r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

      // Authentication routes
      authRouter := r.PathPrefix("/auth").Subrouter()
      authRouter.HandleFunc("/register", handlers.Register).Methods("POST")
      authRouter.HandleFunc("/login", handlers.Login).Methods("POST")

      // API routes
      apiRouter := r.PathPrefix("/api/v1").Subrouter()
      apiRouter.Use(middleware.JWTAuth)
      apiRouter.Use(middleware.RateLimit)

      // Book management
      apiRouter.HandleFunc("/books", handlers.CreateBook).Methods("POST")
      apiRouter.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
      apiRouter.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
      apiRouter.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")
      apiRouter.HandleFunc("/books", handlers.GetCachedBooks).Methods("GET")
      apiRouter.HandleFunc("/search", handlers.SearchBooksV2).Methods("GET")

      // Uploads
      apiRouter.HandleFunc("/upload", handlers.UploadCoverImage).Methods("POST")

      // Analytics
      apiRouter.HandleFunc("/analytics", handlers.GetAnalytics).Methods("GET")

      port := os.Getenv("PORT")
      if port == "" {
        port = "8080"
      }

      log.Printf("Server running on port %s", port)
      log.Fatal(http.ListenAndServe(":"+port, r))
    }
