package redis

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "log"
)

var Rdb *redis.Client
var Ctx = context.Background()

func Connect(host, password string, port int) {
    addr := fmt.Sprintf("%s:%d", host, port) // формируем адрес с портом
    Rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })

    _, err := Rdb.Ping(Ctx).Result()
    if err != nil {
        log.Fatal("Redis connection failed:", err)
    }

    log.Println("Redis connected:", addr)
}
