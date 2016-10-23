package main

import(
  "fmt"
  "net/http"
  "strconv"
  "encoding/json"

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

func Status(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "nice to see you")
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  redisCounter, err := redis.HGetAll(params["counterName"]).Result()
  LogError(err)

  if len(redisCounter) == 0{
    fmt.Fprintln(w, "{message: counter",  params["counterName"] ,"not found}",)
    return
  }

  counter, createError := model.CreateCounterFromMap(params["counterName"], redisCounter)

  LogError(createError)

  fmt.Fprintln(w, counter.ToJson())
}

func PopulateCounters(w http.ResponseWriter, r *http.Request) {
  counterSize, err := redis.LLen(globalCounterKeys).Result()
  LogError(err)

  if counterSize == 0 {
    fmt.Fprintln(w, "{message: Your counter is empty. add some first}")
    return
  }

  var counters model.Counters
  redisResult, err := redis.LRange(globalCounterKeys, 0, counterSize).Result()
  LogError(err)

  for _, counterKey := range redisResult {
    redisValue, _ := redis.HGet(counterKey, "value").Result()
    redisInitialValue, _ := redis.HGet(counterKey, "initial_value").Result()

    value, _ := strconv.Atoi(redisValue)
    initialValue, _ := strconv.Atoi(redisInitialValue)

    counter := model.Counter{Name: counterKey, Value: value, InitialValue: initialValue}

    counters = append(counters, &counter)
  }
  json.NewEncoder(w).Encode(counters)
}

func ValidateParams(params map[string]string) (bool){
  status := true

  if params["name"] == "" {
    status = false
  } else if params["initial_value"] == "" {
    status = false
  }

  if len(params) == 0 {
    status = false
  }

  return status
}

func NewCounter(w http.ResponseWriter, r *http.Request) (){
  params :=  make(map[string]string)
  params["name"] = r.FormValue("name")
  params["initial_value"] = r.FormValue("initial_value")

  fmt.Println("params[initial_value]", params["initial_value"])
  fmt.Println("params[name] ", params["name"])
  if ValidateParams(params) == false {
    fmt.Fprintln(w, "{ message: Invalid parameter }")
    return
  }

  initial_value, _ := strconv.Atoi(params["initial_value"])
  counter, err := model.AddNew(params["name"], initial_value, redis)
  LogError(err)

  if err != nil {
    fmt.Fprintln(w, "{message: Error creating counter}")
    return
  }

  if counter != nil {
    redis.LPush(globalCounterKeys, counter.Name)
  }

  fmt.Fprintln(w, counter.ToJson())
}

func main() {
  gotenv.Load()

  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/status", Status)
  router.HandleFunc("/counter/all", PopulateCounters)
  router.HandleFunc("/counter/{counterName}", GetCounter).Methods("GET")
  router.HandleFunc("/counter/new", NewCounter).Methods("POST")

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
