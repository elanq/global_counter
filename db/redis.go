package db

import(
  "fmt"
  "gopkg.in/redis.v5"
  "os"
)

var redisClient *redis.Client

func ConnectRedis() (*redis.Client, error) {
  var err error
  redisClient = redis.NewClient(&redis.Options {
    Addr: os.GetEnv("REDIS_HOST"),
    })

  _, err = redisClient.Ping().Result()

  if err != nil {
    fmt.Println("Failed to connect redis")
    fmt.Println("caused by: ", err)
  }
  return redisClient, err
}
