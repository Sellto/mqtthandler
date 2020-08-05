package lib


// Struct that be used to parse the Message
type Message struct {
  Action string
  Payload interface{}
}

type RunPayload struct {
  Launcher string
  Args []string
}
