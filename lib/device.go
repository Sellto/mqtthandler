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
)

type Device struct {
  Model string `yaml:"model"`
  Launcher string `yaml:"launcher"`
  MQTT MQTT `yaml:"broker"`
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


func (d *Device) action(m Message, client mqtt.Client) {
  switch m.Action {
    case "Ping":
  		d.pong(m,client)
  	case "Info" :
  		d.info(client)
    case "Download" :
      d.download(m,client)
  	}
}


func (d *Device) pong(m Message, client mqtt.Client) {
  if m.Payload == d.MQTT.ID{
    log.Println("Ping Requested")
    if token := client.Publish("response", 0, false, "Pong"); token.Wait() && token.Error() != nil {
      log.Fatal(token.Error())
    }
    log.Println("Pong successfully sent")
  }
}

func (d *Device) info(client mqtt.Client) {
  log.Println("Information of the device requested")
  b := fmt.Sprintf(`"ID":"%s","Model":"%s","Launcher":"%s"}`,d.MQTT.ID,d.Model,d.Launcher)
  if token := client.Publish("response", 0, false, b); token.Wait() && token.Error() != nil {
    log.Fatal(token.Error())
  }
  log.Println("Information successfully sent")
}

func (d *Device) download(m Message, client mqtt.Client) {
  log.Println("Download Request")
  out, err := exec.Command("wget",m.Payload).Output()
	if err != nil {
		log.Println(err)
	}
  if token := client.Publish("response", 0, false, string(out)); token.Wait() && token.Error() != nil {
    log.Fatal(token.Error())
  }
  log.Println("Downloaded at")
}
