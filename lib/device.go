package lib

import(
  "log"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os/signal"
  "syscall"
  "os"
  "fmt"
  "os/exec"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  "strings"
)

type Device struct {
  Type string `yaml:"type"`
  Model string `yaml:"model"`
  Launchers []string `yaml:"launchers"`
  MQTT MQTT `yaml:"broker"`
  Path string `yaml:"sourcepath"`
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
  d.MQTT.Action = d.action
  d.MQTT.Connect()
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
  if channel == d.MQTT.Topic{
    switch m.Action {
      case "Ping":
    		d.pong(m,client)
    	case "Info" :
    		d.info(client)
      }
   }

   if channel == d.MQTT.ID+"-incoming"{
     switch m.Action {
       case "Download" :
         d.download(m,client)
       case "Run" :
         d.run(m,client)
       }
    }
}


func (d *Device) info(client mqtt.Client) {
  log.Println("Information of the device requested")
  b := fmt.Sprintf(`"ID":"%s","Model":"%s","Launchers":"%s","Path":"%s"}`,d.MQTT.ID,d.Model,d.Launchers,d.Path)
  //Send the information to the info channel
  if token := client.Publish("info", 0, false, b); token.Wait() && token.Error() != nil {
    log.Fatal(token.Error())
  }
  log.Println("Information successfully sent")
}

func (d *Device) pong(m Message, client mqtt.Client) {
  if m.Payload == d.MQTT.ID{
    log.Println("Ping Requested")
    //Send a Pong into the specific channel
    if token := client.Publish(d.MQTT.ID, 0, false, "Pong"); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Println("Pong successfully sent")
  }
}


func (d *Device) download(m Message, client mqtt.Client) {
  //Parse Payload as a simple string
  s, ok := m.Payload.(string)
  if ok {
    log.Printf("Download of %s request",s)
    //try to download the file with wget tools
    out, err := exec.Command("wget","-P",d.Path,s).Output()
  	if err != nil {
  		log.Println(err)
  	}
    // Send a response
    if token := client.Publish(d.MQTT.ID, 0, false, string("Downloaded at "+d.Path); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Printf("Downloaded at %s",d.Path)
  }
}


func (d *Device) run(m Message, client mqtt.Client) {
  log.Printf("Run application request")
  var p RunPayload
  v, ok := m.Payload.(map[string]interface{})
  if ok {
    p.Launcher, ok = v["launcher"].(string)
    if ok {
      s, ok := v["args"].(string)
      p.Args = strings.Split(s," ")
      if ok {
          //Check if the launcher is available
          if contains(d.Launchers,p.Launcher){
            log.Printf("Launcher %s available on this platform",p.Launcher)
            // Run the requested action
            out, err := exec.Command(p.Launcher,p.Args...).Output()
          	if err != nil {
          		log.Println(err)
          	}
            //Send result to MQTT
            if token := client.Publish(d.MQTT.ID, 0, false, string(out)); token.Wait() && token.Error() != nil {
              log.Fatal(token.Error())
            }
          } else {
            log.Println("Launcher not available on this platform")
          }
        } else {
            log.Println("can't parse the args field")
        }
      } else {
        log.Println("can't parse the launcher field")
      }
    } else {
      log.Println("can't parse the payload field")
    }
  }
