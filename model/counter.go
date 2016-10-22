package model

import(
  "fmt"
  "gopkg.in/redis.v5"
  "strconv"
  "encoding/json"
)

type Counter struct {
  Name          string  `json:name`
  Value         int     `json:value`
  InitialValue  int     `json:initial_value`
}

func AddNew(counterName string, initialValue int, redis *redis.Client) (*Counter, error) {
  counter := &Counter{Name: counterName, Value: initialValue, InitialValue: initialValue}
  return counter.UpdateValue(initialValue, 0, redis)
}

func (counter *Counter) UpdateValue(value int, operation int, redis *redis.Client) (*Counter, error){

  var err error

  switch operation {
  case 0:
    counterValue := strconv.Itoa(counter.Value)
    initialValue := strconv.Itoa(counter.InitialValue)
    redis.HSet(counter.Name, "value", counterValue)
    redis.HSet(counter.Name, "initial_value", initialValue)
  default:
    err = fmt.Errorf("Unrecognize update operation")
  }

  return counter, err
}

func (counter *Counter) ToJson() (string) {
  json, _ := json.Marshal(counter)
  return string(json)
}

