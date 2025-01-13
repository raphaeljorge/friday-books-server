# Book Management API

    ## New Features
    - Background task processing with Redis
    - Advanced search with faceted filtering
    - Full-text search capabilities
    - Comprehensive health check endpoint
    - Security headers middleware
    - CORS support
    - Search index implementation
    - Background indexing of books

    ## Usage
    ### Background Tasks
    To enqueue a background task:
    ```go
    err := handlers.EnqueueTask("book:index", book)
    if err != nil {
      log.Printf("Error enqueueing task: %v", err)
    }
    ```

    ### Advanced Search
    Example search request:
    ```bash
    curl -X GET "http://localhost:8080/api/v1/search?q=harry+potter&genre=fantasy&language=en&format=hardcover"
    ```

    ### Health Check
    ```bash
    curl http://localhost:8080/health
    ```

    ## Deployment
    Make sure to run the background worker:
    ```bash
    ./book-api -worker
    ```
