package main

import(
  lib "./lib"
)

func main(){
  server := lib.Server{}
  server.GetConfigFromFile("config.yaml")
  //device.MQTT.ID="mac"
  server.Listen()
}
