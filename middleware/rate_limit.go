package middleware

    import (
      "net/http"
      "github.com/ulule/limiter/v3"
      "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
      "github.com/ulule/limiter/v3/drivers/store/memory"
    )

    func RateLimit(next http.Handler) http.Handler {
      rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
      }

      store := memory.NewStore()
      instance := limiter.New(store, rate)
      middleware := stdlib.NewMiddleware(instance)

      return middleware.Handler(next)
    }
