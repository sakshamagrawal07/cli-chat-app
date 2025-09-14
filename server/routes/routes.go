package routes

import (
	"net/http"

	"github.com/sakshamagrawal07/cli-chat-app.git/server/handler"
)

func SetupRoutes() {
	http.HandleFunc("/ws", handler.HandleWebSocket)
}