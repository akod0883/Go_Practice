package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clients       = make(map[*websocket.Conn]bool)
	clientsMux    sync.Mutex
	messageBuffer = make(chan []byte, 10)
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("public")))
	log.Println("Server started at http://localhost:8080")
	go processMessages()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	clientsMux.Lock()
	clients[conn] = true
	clientsMux.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}

		messageBuffer <- message
	}

	clientsMux.Lock()
	delete(clients, conn)
	clientsMux.Unlock()
}

func processMessages() {
	for {
		message := <-messageBuffer
		clientsMux.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error broadcasting message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMux.Unlock()
		fmt.Println(string(message))
	}
}
