package chat

type Hub struct {
	connections map[*Connection]bool
	broadcast chan []byte
	register chan *Connection
	unregister chan *Connection
}
 
var H = Hub {
	broadcast:   make(chan []byte),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
}

func (this *Hub) Run() {
	for {
		select {
			case c := <- this.register:
				this.connections[c] = true
			case c := <- this.unregister:
				if _, ok := this.connections[c]; ok {
					delete(this.connections, c)
					close(c.send)
				}
			case msg := <- this.broadcast:
				for c := range this.connections {
					select {
						case c.send <- msg:
						default:
							delete(this.connections, c)
							close(c.send)
					}
				}
		}
	}
}