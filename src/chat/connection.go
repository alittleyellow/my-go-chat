package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Connection struct {
	ws *websocket.Conn
	send chan []byte
}

func (this *Connection) reader() {
	for {
		_, message, err := this.ws.ReadMessage()
		if err != nil {
			break
		}
		H.broadcast <- message
	}
	this.ws.Close()
}

func (this *Connection) writer() {
	for message := range this.send {
		err := this.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	this.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return 
	}
	conn := &Connection{send: make(chan []byte, 256), ws: ws}
	H.register <- conn
	go conn.writer()
	conn.reader()
	defer func() {H.unregister <- conn }()
}








