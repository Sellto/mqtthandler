package main

import(
  mqtt "github.com/eclipse/paho.mqtt.golang"
  "fmt"
  "log"
  "os"
  "os/exec"
  "os/signal"
  "syscall"
  "io/ioutil"
  lib "./lib"
)


func ParseConfigFile(filename string)lib.Config{
  conf := lib.Config{}
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    fmt.Println(err)
  }
  conf.GetAllFromFile(yamlFile)
  return conf
}

func MqttHandler(conf lib.Config)func(client mqtt.Client, message mqtt.Message) {
var msgRcvd mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
  for _,sub := range conf.Subscribers {
    if sub.Topic == message.Topic() {
      for _,trigger := range conf.Triggers {
        if trigger.Subscriber == sub.Name && trigger.TriggerWord == string(message.Payload()){
          for _,action := range conf.Actions {
            if trigger.Actions == action.Name {
              for _,a := range action.ShellActions {
                out, err := exec.Command(a.Cmd,a.Arg).Output()
                	if err != nil {
                		log.Fatal(err)
                	}
                  if token := client.Publish("response", 0, false, string(out)); token.Wait() && token.Error() != nil {
                     log.Fatal(token.Error())
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  return msgRcvd
}


func SetupCloseHandler(clients []mqtt.Client) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
    for _,c := range clients {
      c.Disconnect(250)
    }
		os.Exit(0)
	}()
}

func MqttConnect(conf lib.Config)[]mqtt.Client{
  var clients []mqtt.Client
  for _, sub := range conf.Subscribers {
    for _, broker := range conf.Brokers {
      if sub.Broker == broker.Name {
        opts := mqtt.NewClientOptions().AddBroker(broker.Host+":"+broker.Port).SetClientID(sub.Name)
        c := mqtt.NewClient(opts)
        if token := c.Connect(); token.Wait() && token.Error() != nil {
           log.Fatal(token.Error())
        }
        if token := c.Subscribe(sub.Topic, 0, MqttHandler(conf)); token.Wait() && token.Error() != nil {
           log.Fatal(token.Error())
        }
        clients = append(clients, c)
        }
      }
  }
  return clients
}



func main(){
  SetupCloseHandler(MqttConnect(ParseConfigFile("config.yaml")))
  for {
  }
}
