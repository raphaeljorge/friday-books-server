package storage

    import (
      "context"
      "github.com/go-redis/redis/v8"
      "os"
      "time"
    )

    var ctx = context.Background()

    func GetRedis() *redis.Client {
      return redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
      })
    }

    func CacheSet(key string, value interface{}, expiration time.Duration) error {
      return GetRedis().Set(ctx, key, value, expiration).Err()
    }

    func CacheGet(key string) (string, error) {
      return GetRedis().Get(ctx, key).Result()
    }

    func CacheDelete(key string) error {
      return GetRedis().Del(ctx, key).Err()
    }
