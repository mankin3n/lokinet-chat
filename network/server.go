package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Chatroom struct {
	Name    string
	Clients map[net.Conn]string
	Mutex   sync.Mutex
}

var (
	chatrooms   = make(map[string]*Chatroom) // List of chatrooms
	chatroomsMu sync.Mutex                   // Mutex to protect chatrooms map
)

// StartServer initializes the chatroom server.
func StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Chatroom server started on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		go handleClient(conn)
	}
}

// handleClient manages communication with a single client.
func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("Enter the chatroom name: "))
	chatroomName, _ := reader.ReadString('\n')
	chatroomName = strings.TrimSpace(chatroomName)

	chatroom := getOrCreateChatroom(chatroomName)

	conn.Write([]byte("Enter your username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	chatroom.Mutex.Lock()
	chatroom.Clients[conn] = username
	chatroom.Mutex.Unlock()

	broadcast(chatroom, fmt.Sprintf("%s has joined the chatroom\n", username), conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			chatroom.Mutex.Lock()
			delete(chatroom.Clients, conn)
			chatroom.Mutex.Unlock()
			broadcast(chatroom, fmt.Sprintf("%s has left the chatroom\n", username), conn)
			return
		}

		broadcast(chatroom, fmt.Sprintf("%s: %s", username, message), conn)
	}
}

// getOrCreateChatroom retrieves an existing chatroom or creates a new one.
func getOrCreateChatroom(name string) *Chatroom {
	chatroomsMu.Lock()
	defer chatroomsMu.Unlock()

	if chatroom, exists := chatrooms[name]; exists {
		return chatroom
	}

	newChatroom := &Chatroom{
		Name:    name,
		Clients: make(map[net.Conn]string),
	}
	chatrooms[name] = newChatroom
	return newChatroom
}

// broadcast sends a message to all connected clients in the chatroom except the sender.
func broadcast(chatroom *Chatroom, message string, sender net.Conn) {
	chatroom.Mutex.Lock()
	defer chatroom.Mutex.Unlock()

	fmt.Printf("[%s] Broadcasting: %s", chatroom.Name, message)
	for client := range chatroom.Clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}
