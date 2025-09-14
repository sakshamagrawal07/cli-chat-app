package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	kafkapkg "github.com/sakshamagrawal07/cli-chat-app.git/server/kafka"
	redispkg "github.com/sakshamagrawal07/cli-chat-app.git/server/redis"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

// SafeConn wraps a websocket.Conn with a mutex
type SafeConn struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

func (sc *SafeConn) WriteJSON(v interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.Conn.WriteJSON(v)
}

// Global clients map
var Clients = make(map[*SafeConn]string)
var subscribedChannels = make(map[string]bool) // tracks which channels we’ve already subscribed to
var mu sync.Mutex                             // protects Clients + subscribedChannels

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}
	defer conn.Close()

	safeConn := &SafeConn{Conn: conn}
	mu.Lock()
	Clients[safeConn] = ""
	mu.Unlock()

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Disconnected:", Clients[safeConn])
			mu.Lock()
			delete(Clients, safeConn)
			mu.Unlock()
			redispkg.MarkUserOffline(msg.SenderUsername)
			return
		}

		if Clients[safeConn] == "" {
			mu.Lock()
			Clients[safeConn] = msg.SenderUsername
			mu.Unlock()
			redispkg.MarkUserOnline(msg.SenderUsername)
		}

		log.Printf("Received from %s: %s\n", msg.SenderUsername, msg.Message)

		// ✅ Subscribe only once per channel
		mu.Lock()
		if !subscribedChannels[msg.Channel] {
			subscribedChannels[msg.Channel] = true
			go redispkg.Subscribe(Broadcast, msg.Channel) // run subscription in background goroutine
		}
		mu.Unlock()	

		// Send to Kafka for persistence
		if err := kafkapkg.ProduceMessage(msg); err != nil {
			log.Println("Kafka error:", err)
		}

		// Publish via Redis (no duplicate subscribes!)
		redispkg.Publish(msg)
	}
}

func Broadcast(msg models.Message) {
	mu.Lock()
	defer mu.Unlock()

	for conn, username := range Clients {
		if username == msg.RecipientUsername || username == msg.SenderUsername {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Broadcast error:", err)
			}
		}
	}
}
