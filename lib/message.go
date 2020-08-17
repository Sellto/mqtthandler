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


// Info request
// send to global Topic, response on "info" channel
//
// {
// 	"action":"Info",
// 	"payload":""
// }

// Response :
// {"ID":"crimson-flower","Model":"MacBook-Pro","Launchers":"[python3 go]","Path":"/Users/tse/data/mqtthandler/"}



// send to global Topic, response on id_of_the_device channel
//
// {
// 	"action":"Ping",
// 	"payload": "id-of-the-device"
// }

// Response :
// Pong

// send to id_of_the_device+-incoming, response on id_of_the_device channel
//
// {
// 	"action":"Download",
// 	"payload": "url-of-the-file"
// }

// Response :
// Downloaded at {{ sourcepath }}
