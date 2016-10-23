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

func CreateCounterFromMap(key string, counterMap map[string]string) (Counter, error) {
  var err error
  if counterMap == nil {
    err = fmt.Errorf("map is nil")
  }

  value, _ := strconv.Atoi(counterMap["value"])
  initialValue, _ := strconv.Atoi(counterMap["initial_value"])

  counter := Counter{
    Name:         key,
    Value:        value,
    InitialValue: initialValue}
  return counter, err
}

func AddNew(counterName string, initialValue int, redis *redis.Client) (*Counter, error) {
  counter := &Counter{Name: counterName, Value: initialValue, InitialValue: initialValue}
  return counter.UpdateValue(initialValue, 0, redis)
}

func (counter *Counter) UpdateValue(value int, operation int, redis *redis.Client) (*Counter, error){

  var err error

  switch operation {
  case 0:
    redis.HSet(counter.Name, "value", strconv.Itoa(counter.Value))
    redis.HSet(counter.Name, "initial_value", strconv.Itoa(counter.InitialValue))
  case 1:
    //add
    totalValue := counter.Value + value
    counter.Value = totalValue
    redis.HSet(counter.Name, "value", strconv.Itoa(totalValue))
  case 2:
    //substract
    substractValue := counter.Value - value
    counter.Value = substractValue
    redis.HSet(counter.Name, "value", strconv.Itoa(substractValue))
  default:
    err = fmt.Errorf("Unrecognize update operation")
  }

  return counter, err
}

func (counter *Counter) ToJson() (string) {
  json, _ := json.Marshal(counter)
  return string(json)
}

