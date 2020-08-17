package lib

import (
  "encoding/json"
  //"fmt"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}


func JsonMessage(m string,p interface{}) []byte {
  var message Message
  message.Action = m
  message.Payload = p
  b,_ :=  json.Marshal(message)
  //fmt.Println(string(b))
  return b
}
