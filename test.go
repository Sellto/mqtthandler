package main

import(
  lib "./lib"
  "fmt"
)

func main(){
  str := lib.JsonMessage("test",lib.Device{})
  fmt.Println(str)
}
