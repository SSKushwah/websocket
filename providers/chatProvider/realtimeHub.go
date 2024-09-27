package chatProvider

import (
	"websocket/providers"
	"websocket/utils"
)

type RealtimeHub struct {
	clients map[string]*NewClient

	// Accept registration requests from the clients.
	register chan *NewClient

	// Accept deregistration requests from the clients.
	deregister chan *NewClient

	// Send message to the client with the specified Name.
	sendToClientName chan string
}

func NewRealtimeChatProvider() providers.RealtimeChatHubProvider {

	return &RealtimeHub{
		register:         make(chan *NewClient, 1),
		deregister:       make(chan *NewClient, 1),
		clients:          make(map[string]*NewClient),
		sendToClientName: make(chan string, 1),
	}
}

// Run initiates an infinite for loop for instantly communicate between gRPC service and server via go channels
func (realtimeHub *RealtimeHub) Run() {
	for {
		select {
		case clientClient := <-realtimeHub.register:
			if _, ok := realtimeHub.clients[clientClient.Name]; !ok {
				var newClient *NewClient
				realtimeHub.clients[clientClient.Name] = newClient
			}
			realtimeHub.clients[clientClient.Name] = clientClient

		case clientClient := <-realtimeHub.deregister:
			delete(realtimeHub.clients, clientClient.Name)
		case sendToclient := <-realtimeHub.sendToClientName:
			utils.LogDebug("Run", "HUB is running an infnitely to add clients and communicate with them", "", sendToclient)
		}
	}
}

// Get gets the pointer client struct
func (RealtimeHub *RealtimeHub) Get() interface{} {
	return RealtimeHub
}

// Stop stops the send channel
func (RealtimeHub *RealtimeHub) Stop() {
	for _, userClients := range RealtimeHub.clients {
		close(userClients.Send)
	}
}
