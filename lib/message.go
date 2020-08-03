package lib


// Struct that be used to parse the Message
type Message struct {
  Action string
  Payload string
}

// Accepted message
// {
//   "action":"a specific action",
//   "payload:"some text"
// }
