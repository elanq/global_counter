package model

import(
  "encoding/json"
  "github.com/elanq/global_counter/model"
)

type Response struct {
  Status string `json:status`
  Message string `json:message`
}

type CounterResponse struct {
  Response *Response
  Counter model.Counter
}

func (response *Response) ToJson() (string) {
  json, _ := json.Marshal(response)
  return string(json)
}
