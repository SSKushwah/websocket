package chatProvider

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
	"websocket/models"
	"websocket/utils"

	"github.com/gorilla/websocket"
)

type NewClient struct {
	Name         string
	Hostname     string
	IPAddress    string
	Platform     string
	OSAndVersion string
	HUB          *RealtimeHub
	Conn         *websocket.Conn
	Send         chan models.SendMessage
	Timer        time.Timer
	Ticker       *time.Ticker
}

// NewClientStream gets the metaData and stream of agent and keep it in the agent HUB map.
// func (s *websocket) NewClientStream(hub *RealtimeHub, clientContext *models.ClientContext, realtimeHub providers.RealtimeChatHubProvider) *NewClient {
func NewClientStream(hub *RealtimeHub, clientContext *models.ClientContext, conn *websocket.Conn) *NewClient {

	ticker := time.NewTicker(10 * time.Second)
	return &NewClient{
		Name:     clientContext.Name,
		Hostname: clientContext.Hostname,
		Platform: clientContext.Platform,
		HUB:      hub,
		Send:     make(chan models.SendMessage, 1),
		Timer:    time.Timer{},
		Ticker:   ticker,
	}
}

func (newClient *NewClient) Get() *NewClient {
	return newClient
}

func (newClient *NewClient) Register() {
	newClient.HUB.register <- newClient
}

func (newClient *NewClient) Unregister() {
	newClient.HUB.deregister <- newClient
}

func (n *NewClient) WritePump() {
	count := 0
	for {
		fmt.Println("iteration", count)

		count++
		select {

		// wait till send channel is empty. after the write operation on send channel stream sends the message to server.
		case sendMessage, ok := <-n.Send:
			if ok {

				err := n.Conn.WriteMessage(1, sendMessage.Message)

				if err != nil {
					utils.LogError("client.go", "WritePump :error sending messages to the client", n.Name, err)
				}
				utils.LogInfo("Write", "sent the message successfully", sendMessage.MessageType, nil)
			}
		case <-n.Timer.C:
			utils.LogInfo("client.go", "WritePump :ping timer finished", n.Name, nil)
			return
		}
	}
}

func (n *NewClient) ReadPump() {
	count := 0
	for {
		count++

		_, clientMsg, err := n.Conn.ReadMessage()
		// clientMsg, err := n.Stream.Recv()
		if err != nil {
			if err == io.EOF {
				continue
			}
			utils.LogError("client.go", "ReadPump: error getting message from server stream.", "", err)
			break
		} else {

			utils.LogInfo("Read", "read the message successfully", string(clientMsg), nil)

			// comment the log below if dont want to see the messages on console
			log.Printf("Receving message : %s\n", string(clientMsg))

			// go n.ProcessClientMessaging(clientMsg.MessageType, clientMsg)

		}
	}
}

func (n *NewClient) ProcessClientMessaging(messageType string, message []byte) {

	switch messageType {

	case models.PingMessageType:
		go n.ProcessPing()

	default:
		err := errors.New("invalid message type")
		utils.LogWarning("ProcessClientMessaging", err.Error(), messageType, err)

	}

}

func (nc *NewClient) ProcessPing() {
	log.Printf("Sending  message : %s\n", models.PongMessage)
	nc.Timer = time.Timer{C: time.After(10 * time.Second)}
	nc.Send <- models.SendMessage{
		Message:     []byte(models.PongMessage),
		MessageType: models.PongMessageType,
	}
}
