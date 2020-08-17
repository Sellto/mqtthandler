package lib

import(
  "log"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os/signal"
  "syscall"
  "os"
  "fmt"
  //"os/exec"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  //"strings"
  //"encoding/json"
  "time"
  "github.com/mitchellh/mapstructure"
  "crypto/rand"
)

type Server struct {
  MQTT MQTT `yaml:"broker"`
  Devices map[string]Device
  Token string `yaml:"-"`
}


func (s *Server) GetConfigFromFile(filename string){
  // Read file
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Println(err)
  }
  // Parse file
  err = yaml.Unmarshal(yamlFile, &s)
  if err != nil {
     log.Fatal(err)
  }
}


func (s *Server) Listen() {
  s.Devices = map[string]Device{}
  b := make([]byte,12)
  rand.Read(b)
  s.Token = fmt.Sprintf("%x",b)
	c := make(chan os.Signal)
  s.MQTT.Action = s.action
  s.MQTT.Connect()
  //Signal that trig a Keyboard Interrupt
  s.MQTT.Subscribe("bootstrap")
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  // Goroutine that handle the Keyboard Interrupt
	go func() {
		<-c
		log.Println("\r- Ctrl+C pressed in Terminal")
    // Disconnect to the MQTT Broker
    s.MQTT.Client.Disconnect(250)
    // Close Application
		os.Exit(0)
	}()
  // Still App Open
  go s.Heartbeat(10)
  log.Printf("Server is reachable with the token %s",s.Token)
  for{
  }
}

func (s *Server) Heartbeat(sec int) {
  for {
    time.Sleep(time.Duration(sec)*time.Second)
    for _,d := range s.Devices {
      if token := s.MQTT.Client.Publish(d.ID+"-incoming", 0, false, JsonMessage("Update",d.ID)); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
      }
    }
  }
}


func (s *Server) action(m Message, client mqtt.Client,channel string) {
  if channel == "bootstrap" {
    var device Device
    mapstructure.Decode(m.Payload, &device)
    if device.Token == s.Token {
      log.Printf("Token match")
      s.Devices[device.ID] = device
      log.Printf("New device registred with ID : %s", device.ID)
    s.MQTT.Subscribe(device.ID)
    if token := s.MQTT.Client.Publish(device.ID+"-incoming", 0, false, JsonMessage("Register","Registered")); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
  }
   } else if _, ok := s.Devices[channel]; ok  {
     var device Device
     mapstructure.Decode(m.Payload, &device)
     log.Printf("Device %s updated",device.ID)
     s.Devices[channel] = device
     //log.Println(device)
   }
}
