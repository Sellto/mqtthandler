package lib

import(
  "log"
  "gopkg.in/yaml.v2"
)

type Broker struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
  Host string `yaml:"host"`
  Port string `yaml:"port"`
}

type Subscriber struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
  Broker string `yaml:"broker"`
  Topic string `yaml:"topic"`
}

type ShellAction struct {
  Cmd string `yaml:"cmd"`
  Arg string `yaml:"arg"`
}

type Action struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
  ShellActions []ShellAction `yaml:"shell"`
}

type Trigger struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
  Subscriber string `yaml:"subscriber"`
  TriggerWord string `yaml:"triggerword"`
  Actions string `yaml:"actions"`
}

type Config struct {
  Brokers []Broker
  Subscribers []Subscriber
  Actions []Action
  Triggers []Trigger
}

func (c *Config) getBrokersFromFile(yamlFile []byte){
  input := []Broker{}
  err := yaml.Unmarshal(yamlFile, &input)
  if err != nil {
     log.Fatal(err)
  }
  for _,v := range(input){
      if v.Type == "broker" {
        c.Brokers = append(c.Brokers, v)
      }
    }
}

func (c *Config) getSubscribersFromFile(yamlFile []byte){
  input := []Subscriber{}
  err := yaml.Unmarshal(yamlFile, &input)
  if err != nil {
     log.Fatal(err)
  }
  for _,v := range(input){
      if v.Type == "mqtt-subscriber" {
        c.Subscribers = append(c.Subscribers, v)
      }
    }
}

func (c *Config) getActionsFromFile(yamlFile []byte){
  input := []Action{}
  err := yaml.Unmarshal(yamlFile, &input)
  if err != nil {
     log.Fatal(err)
  }
  for _,v := range(input){
      if v.Type == "action" {
        c.Actions = append(c.Actions, v)
      }
    }
}

func (c *Config) getTriggersFromFile(yamlFile []byte){
  input := []Trigger{}
  err := yaml.Unmarshal(yamlFile, &input)
  if err != nil {
     log.Fatal(err)
  }
  for _,v := range(input){
      if v.Type == "trigger" {
        c.Triggers = append(c.Triggers, v)
      }
    }
}

func (c *Config) GetAllFromFile(yamlFile []byte){
  c.getBrokersFromFile(yamlFile)
  c.getSubscribersFromFile(yamlFile)
  c.getActionsFromFile(yamlFile)
  c.getTriggersFromFile(yamlFile)
}
