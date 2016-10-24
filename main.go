package main

import(
  "fmt"
  "net/http"
  "strconv"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/subosito/gotenv"

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
  json.NewEncoder(w).Encode(redisCounter)
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

func ValidateParams(params map[string]string, requiredArgs int) (bool){
  if len(params) == 0 || len(params) < requiredArgs {
    return false
  }

  return true
}

func NewCounter(w http.ResponseWriter, r *http.Request) (){
  params :=  make(map[string]string)
  params["name"] = r.FormValue("name")
  params["initial_value"] = r.FormValue("initial_value")

  fmt.Println("params[initial_value]", params["initial_value"])
  fmt.Println("params[name] ", params["name"])
  if ValidateParams(params, 2) == false {
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

  json.NewEncoder(w).Encode(counter)
}

func UpdateCounterValue(w http.ResponseWriter, r *http.Request) {
  var err error
  params := make(map[string]string)
  params["name"] = r.FormValue("name")
  params["value"] = r.FormValue("value")
  params["operation"] = r.FormValue("operation")

  if ValidateParams(params, 3) == false {
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

  json.NewEncoder(w).Encode(updatedCounter)

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

  // fmt.Println("Testing counter")
  // AddCounter("beats_audio", 2000)
  // AddCounter("tissue_paseo", 3000)
  // fmt.Println("Populating counter")
  // for _, counter := range counters {
  //   fmt.Println(counter.ToJson())
  // }
}
