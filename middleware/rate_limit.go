package middleware

    import (
      "net/http"
      "time"

      "github.com/ulule/limiter/v3"
      "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
      "github.com/ulule/limiter/v3/drivers/store/memory"
    )

    // Rest of the rate_limit.go file remains the same
