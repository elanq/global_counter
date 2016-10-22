package main

import(
  "fmt"
  "global-counter/db"
  "global-counter/model"
)

var(
  redis, redisErr = db.ConnectRedis()
  counters []*model.Counter
)


func LogError(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

func AddCounter(counterName string, initialValue int) {
  newCounter, err := model.AddNew(counterName, initialValue, redis)
  LogError(err)
  fmt.Println("Counter %s created", newCounter.Name)
  counters = append(counters, newCounter)
}

func main() {
  fmt.Println("Testing counter")
  AddCounter("beats_audio", 2000)
  AddCounter("tissue_paseo", 3000)
  fmt.Println("Populating counter")
  for _, counter := range counters {
    fmt.Println(counter.ToJson())
  }
}
