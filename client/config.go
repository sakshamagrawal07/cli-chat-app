package main

var (
	ServerHost = "localhost"
	ServerPort = "8080"
	Username   = "anonymous"
)

func GetWebSocketURL() string {
	return "ws://" + ServerHost + ":" + ServerPort + "/ws"
}
