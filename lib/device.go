package lib

import(
  "log"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os/signal"
  "syscall"
  "os"
  //"fmt"
  "os/exec"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  //"encoding/json"
  "strings"
  "time"
)

type Device struct {
  Type string `yaml:"type" json:"type"`
  Model string `yaml:"model" json:"model"`
  Launchers []string `yaml:"launchers" json:"launchers"`
  MQTT MQTT `yaml:"broker" json:"-"`
  Path string `yaml:"sourcepath" json:"path"`
  Ready bool `json:"ready"`
  Registred bool `json:"-"`
  ID string `json:"ID"`
  Token string
}


func (d *Device) GetConfigFromFile(filename string){
  // Read file
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Println(err)
  }
  // Parse file
  err = yaml.Unmarshal(yamlFile, &d)
  if err != nil {
     log.Fatal(err)
  }
}


func (d *Device) Listen() {
	c := make(chan os.Signal)
  d.Ready = true
  d.MQTT.Action = d.action
  d.MQTT.Connect()
  d.MQTT.Subscribe(d.MQTT.Topic)
  d.MQTT.Subscribe(d.MQTT.ID+"-incoming")
  d.ID = d.MQTT.ID
  for !d.Registred {
    d.register(Message{Action:"Registering"},d.MQTT.Client)
    time.Sleep(2*time.Second)
  }
  log.Println("Device Registred")
  //Signal that trig a Keyboard Interrupt
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  // Goroutine that handle the Keyboard Interrupt
	go func() {
		<-c
		log.Println("\r- Ctrl+C pressed in Terminal")
    // Disconnect to the MQTT Broker
    d.MQTT.Client.Disconnect(250)
    // Close Application
		os.Exit(0)
	}()
  // Still App Open

  for{}
}


func (d *Device) action(m Message, client mqtt.Client,channel string) {
  // if channel == d.MQTT.Topic{
  //   switch m.Action {
  //     case "Ping":
  //   		d.pong(m,client)
  //   	//case "Info" :
  //   		//d.info(client)
  //     }
  //  }

   if channel == d.MQTT.ID+"-incoming"{
     switch m.Action {
       case "Register":
         d.register(m,client)
       case "Update":
         d.update(m,client)
       case "Download" :
         d.download(m,client)
       case "Run" :
         d.run(m,client)
       }
    }
}


func (d *Device) register(m Message,client mqtt.Client) {
  if m.Action == "Registering" {
    log.Println("Registering Request send ...")
    // var message Message
    // message.Action ="Registering"
    // message.Payload = d
    // b,_ :=  json.Marshal(message)
    //Send the information to the info channel
    if token := client.Publish("bootstrap", 0, false,JsonMessage("Registering",d)); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Println("Information successfully sent")
  } else if m.Action == "Register" {
    if m.Payload == "Registered" {
      d.Registred = true
    }
  }
}


func (d *Device) update(m Message, client mqtt.Client) {
  if m.Payload == d.MQTT.ID{
    log.Println("Update Requested")
    //Send a Pong into the specific channel
    // var message Message
    // message.Action ="update"
    // message.Payload = d
    // b,_ :=  json.Marshal(message)
    if token := client.Publish(d.MQTT.ID, 0, false, JsonMessage("Update",d)); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Println("Update successfully sent")
  }
}


func (d *Device) download(m Message, client mqtt.Client) {
  //Parse Payload as a simple string
  if s, ok := m.Payload.(string); ok {
    log.Printf("Download of %s request",s)
    //try to download the file with wget tools
    _, err := exec.Command("wget","-P",d.Path,s).Output()
  	if err != nil {
      if token := client.Publish(d.MQTT.ID, 0, false, string("Download Failed")); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
      }
  		log.Println(err)
  	}
    // Send a response
    if token := client.Publish(d.MQTT.ID, 0, false, string("Downloaded at "+d.Path)); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Printf("Downloaded at %s",d.Path)
  }
}


func (d *Device) run(m Message, client mqtt.Client) {
  log.Printf("Run application request")
  var p RunPayload
  if v, ok := m.Payload.(map[string]interface{}); ok {
    if p.Launcher, ok = v["launcher"].(string); ok {
      if s, ok := v["args"].(string); ok {
          p.Args = strings.Split(s," ")
          //Check if the launcher is available
          if contains(d.Launchers,p.Launcher){
            log.Printf("Launcher %s available on this platform",p.Launcher)
            // Run the requested action
            _, err := exec.Command(p.Launcher,p.Args...).Output()
          	if err != nil {
          		log.Println(err)
          	}
            //Send result to MQTT
            if token := client.Publish(d.MQTT.ID, 0, false, "Application Launched"); token.Wait() && token.Error() != nil {
              log.Fatal(token.Error())
            }
          } else {
            if token := client.Publish(d.MQTT.ID, 0, false,"Launcher not available on this platform"); token.Wait() && token.Error() != nil {
              log.Fatal(token.Error())
            }
            log.Println("Launcher not available on this platform")
          }
        } else {
            if token := client.Publish(d.MQTT.ID, 0, false,"Can't parse the args field"); token.Wait() && token.Error() != nil {
              log.Fatal(token.Error())
            }
            log.Println("Can't parse the args field")
        }
      } else {
        if token := client.Publish(d.MQTT.ID, 0, false,"can't parse the launcher field"); token.Wait() && token.Error() != nil {
          log.Fatal(token.Error())
        }
        log.Println("can't parse the launcher field")
      }
    } else {
      if token := client.Publish(d.MQTT.ID, 0, false,"can't parse the payload field"); token.Wait() && token.Error() != nil {
        log.Fatal(token.Error())
      }
      log.Println("can't parse the payload field")
    }
  }
