package server

import (
	"net/http"
	"websocket/models"
	"websocket/providers/chatProvider"
	"websocket/utils"

	"github.com/gorilla/websocket"
)

// var (
// 	/**
// 	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
// 	*/
// 	websocketUpgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 	}
// )

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (srv *Server) Connect(w http.ResponseWriter, r *http.Request) {
	var clientContext models.ClientContext
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.LogError("Connect", "unable to connect to the websocket upgrader", "", err)
		return
	}

	clinetConn := chatProvider.NewClientStream(srv.RealtimeChatProvider.Get().(*chatProvider.RealtimeHub), &clientContext, conn)

	clinetConn.Register()

	go clinetConn.WritePump()
	clinetConn.ReadPump()

}
