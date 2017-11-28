package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var customers Customers

func main() {
	customers.Init()

	http.HandleFunc("/", serveWs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrader error: %s\n", err.Error())
		return
	}
	defer ws.Close()

	reqUrl, _ := url.Parse(r.RequestURI)
	licence := reqUrl.Query().Get("licence")
	if licence == "" {
		sendToWs(ws, "No licence provided")
	}

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			sendToWs(ws, fmt.Sprintf("Error occurred when reading messages: %s\n", err.Error()))
			return
		}
		if messageType != websocket.TextMessage {
			sendToWs(ws, fmt.Sprintf("Received an unexpected message type: %v\n", messageType))
			return
		}

		msg, err := fromJson(p)
		if err != nil {
			sendToWs(ws, fmt.Sprintf("Error decoding: %s", err.Error()))
			return
		}

		err = customers.Add(licence, msg.token, msg.url, msg.fields, ws)
		if err != nil {
			sendToWs(ws, fmt.Sprintf("Couldn't add new record: %s", err.Error()))
			return
		}

		currentCustomers, err := customers.Get(licence)
		if err != nil {
			sendToWs(ws, fmt.Sprintf("Couldn't get customers: %s", err.Error()))
			return
		}

		sendToWs(ws, fmt.Sprintf("Current customers: %+v", currentCustomers))
	}
}

func sendToWs(ws *websocket.Conn, msg string) {
	ws.WriteMessage(websocket.TextMessage, []byte(msg))
}
