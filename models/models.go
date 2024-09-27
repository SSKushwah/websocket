package models

const (
	PingMessage     string = "ping"
	PingMessageType string = "Ping Message Type"
	PongMessage     string = "pong"
	PongMessageType string = "Pong Message Type"
)

// type User struct {
// 	ID         string `json:"id" bson:"id"`
// 	Name       string `json:"name" bson:"name"`
// 	IpAddress  string `json:"ipAddress" bson:"ipAddress"`
// 	MACAddress string `json:"macAddress" bson:"macAddress"`
// }

type ClientContext struct {
	Platform       string `json:"platform" bson:"platform"`
	ComputerSystem `bson:",inline"`
	ClientID       string `json:"clientID" bson:"clientID"`
	Name           string `json:"name" bson:"name"`
}

// Details of the computer system.
type ComputerSystem struct {
	CurrentLoggedInUser string `json:"currentLoggedInUser" bson:"currentLoggedInUser"`
	Domain              string `json:"domain" bson:"domain"`
	Hostname            string `json:"hostname" bson:"hostname"`
}

// Message strcuture to facilitate communication.
type Message struct {
	Message     []byte `json:"message"`
	MessageType string `json:"messageType"`
}

type SendMessage struct {
	MessageType string `json:"messageType"`
	Message     []byte `json:"message"`
}
