package helper

import(
  "net/http"
  "encoding/json"
)

func Response(w http.ResponseWriter, jsonObject interface{}) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(jsonObject)
}

func ValidateParams(params map[string]string, requiredArgs int) (bool){
  if len(params) == 0 || len(params) < requiredArgs {
    return false
  }
  return true
}
