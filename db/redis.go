package db

import(
  "fmt"
  "gopkg.in/redis.v5"
)

var redisClient *redis.Client

func ConnectRedis() (*redis.Client, error) {
  var err error
  redisClient = redis.NewClient(&redis.Options {
    Addr: "127.0.0.1:6379",
    })

  _, err = redisClient.Ping().Result()

  if err != nil {
    fmt.Println("Failed to connect redis")
    fmt.Println("caused by: ", err)
  }
  return redisClient, err
}
