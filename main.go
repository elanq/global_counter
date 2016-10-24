package main

import(
  "fmt"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"
  "github.com/subosito/gotenv"

  "global-counter/helper"
  "global-counter/db"
  "global-counter/model"
)

var(
  _ = gotenv.Load()
  redis, redisErr = db.ConnectRedis()
)

const (
  globalCounterKeys = "global_counter:key_list"
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
  redisCounter, err := model.GetCounter(params["counterName"], redis)
  LogError(err)

  if redisCounter == nil{
    fmt.Fprintln(w, "{message: counter",  params["counterName"] ,"not found}",)
    return
  }
  fmt.Println(redisCounter)

  helper.Response(w, redisCounter)
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
  helper.Response(w, counters)
}

func NewCounter(w http.ResponseWriter, r *http.Request) (){
  params :=  make(map[string]string)
  params["name"] = r.FormValue("name")
  params["initial_value"] = r.FormValue("initial_value")

  fmt.Println("params[initial_value]", params["initial_value"])
  fmt.Println("params[name] ", params["name"])
  if helper.ValidateParams(params, 2) == false {
    fmt.Fprintln(w, "{ message: Invalid parameter }")
    return
  }

  initial_value, _ := strconv.Atoi(params["initial_value"])
  counter, err := model.AddNewCounter(params["name"], initial_value, redis)
  LogError(err)

  if err != nil {
    fmt.Fprintln(w, "{message: Error creating counter}")
    return
  }

  if counter != nil {
    redis.LPush(globalCounterKeys, counter.Name)
  }

  helper.Response(w, counter)
}

func UpdateCounterValue(w http.ResponseWriter, r *http.Request) {
  var err error
  params := make(map[string]string)
  params["name"] = r.FormValue("name")
  params["value"] = r.FormValue("value")
  params["operation"] = r.FormValue("operation")

  if helper.ValidateParams(params, 3) == false {
    fmt.Fprintln(w, "{\"message\": \"Invalid parameter\"}")
    return
  }

  counter, err := model.GetCounter(params["name"], redis)
  LogError(err)
  if err != nil {
    fmt.Fprintln(w, "{\"message\": \"Error populating counter\"}")
    return
  }

  value, _ := strconv.Atoi(params["value"])
  operation, _ := strconv.Atoi(params["operation"])

  if operation == 0 {
    fmt.Fprintln(w, "{\"message\" : \"Invalid operation\"}")
    return
  }
  updatedCounter, err := counter.UpdateValue(value, operation, redis)
  if err != nil {
    fmt.Fprintln(w, "{\"message\":\",", err, "\"}")
    return
  }

  helper.Response(w, updatedCounter)
}

func main() {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/status", Status)
  router.HandleFunc("/counter/all", PopulateCounters)
  router.HandleFunc("/counter/{counterName}", GetCounter).Methods("GET")
  router.HandleFunc("/counter/new", NewCounter).Methods("POST")
  router.HandleFunc("/counter/update", UpdateCounterValue).Methods("PUT")

  fmt.Println("Serving at port 6123")
  http.ListenAndServe(":6123", router)
}
