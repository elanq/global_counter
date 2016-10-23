package main

import(
  "fmt"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"
  "github.com/subosito/gotenv"

  "github.com/elanq/global_counter/db"
  "github.com/elanq/global_counter/model"
)

var(
  redis, redisErr = db.ConnectRedis()
)

const (
  globalCounterKeys = "global_counter_key_list"
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

  redis.LPush(globalCounterKeys, newCounter.Name)
}

func Status(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "nice to see you")
}

func PopulateCounters(w http.ResponseWriter, r *http.Request) {
  counterSize, err := redis.LLen(globalCounterKeys).Result()
  LogError(err)


  if counterSize == 0 {
    fmt.Fprintln(w, "{message: Your counter is empty. add some first}")
    return
  }

  var counters []string
  redisResult, err := redis.LRange(globalCounterKeys, 0, counterSize).Result()
  LogError(err)
  counterKeys := redisResult
  for _, counterKey := range counterKeys {
    redisValue, _ := redis.HGet(counterKey, "value").Result()
    redisInitialValue, _ := redis.HGet(counterKey, "initial_value").Result()

    value, _ := strconv.Atoi(redisValue)
    initialValue, _ := strconv.Atoi(redisInitialValue)

    counter := model.Counter{Name: counterKey, Value: value, InitialValue: initialValue}

    counters = append(counters, counter.ToJson())
  }

  fmt.Fprintln(w, counters)
}

func main() {
  gotenv.Load()

  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/status", Status)
  router.HandleFunc("/counter/all", PopulateCounters)

  fmt.Println("Serving at port 6123")
  http.ListenAndServe(":6123", router)

  // fmt.Println("Testing counter")
  // AddCounter("beats_audio", 2000)
  // AddCounter("tissue_paseo", 3000)
  // fmt.Println("Populating counter")
  // for _, counter := range counters {
  //   fmt.Println(counter.ToJson())
  // }
}
