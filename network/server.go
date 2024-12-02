package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Chatroom struct {
	Name    string
	Clients map[net.Conn]string
	History *os.File
	Mutex   sync.Mutex
}

var (
	chatrooms   = make(map[string]*Chatroom) // List of chatrooms
	chatroomsMu sync.Mutex                   // Mutex to protect chatrooms map
)

// StartServer initializes the chatroom server.
func StartServer() {
	port := "8080" // Default port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// Get the server's IP address
	hostName, _ := os.Hostname()
	addresses, _ := net.LookupHost(hostName)
	var ipAddress string
	if len(addresses) > 0 {
		ipAddress = addresses[0]
	} else {
		ipAddress = "127.0.0.1"
	}

	// Display connection details
	fmt.Printf("Chatroom server started!\n")
	fmt.Printf("Clients can connect using the following details:\n")
	fmt.Printf("IP Address: %s\n", ipAddress)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Example: %s:%s\n\n", ipAddress, port)

	for {
		conn, err := listener.Accept()
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

	conn.Write([]byte("Available chatrooms:\n"))
	listChatrooms(conn)

	conn.Write([]byte("Enter the chatroom name to join or create: "))
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

	conn.Write([]byte("Type your messages. Type '/quit' to leave the chatroom.\n"))

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			chatroom.Mutex.Lock()
			delete(chatroom.Clients, conn)
			chatroom.Mutex.Unlock()
			broadcast(chatroom, fmt.Sprintf("%s has left the chatroom\n", username), conn)
			return
		}

		message = strings.TrimSpace(message)
		if message == "/quit" {
			chatroom.Mutex.Lock()
			delete(chatroom.Clients, conn)
			chatroom.Mutex.Unlock()
			broadcast(chatroom, fmt.Sprintf("%s has left the chatroom\n", username), conn)
			return
		}

		fullMessage := fmt.Sprintf("%s: %s\n", username, message)
		broadcast(chatroom, fullMessage, conn)
		saveToHistory(chatroom, fullMessage)
	}
}

// getOrCreateChatroom retrieves an existing chatroom or creates a new one.
func getOrCreateChatroom(name string) *Chatroom {
	chatroomsMu.Lock()
	defer chatroomsMu.Unlock()

	if chatroom, exists := chatrooms[name]; exists {
		return chatroom
	}

	// Create new chatroom
	file, err := os.OpenFile(fmt.Sprintf("%s_history.txt", name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating history file for chatroom %s: %v\n", name, err)
	}

	newChatroom := &Chatroom{
		Name:    name,
		Clients: make(map[net.Conn]string),
		History: file,
	}
	chatrooms[name] = newChatroom
	return newChatroom
}

// listChatrooms sends the list of available chatrooms to a client.
func listChatrooms(conn net.Conn) {
	chatroomsMu.Lock()
	defer chatroomsMu.Unlock()

	for name := range chatrooms {
		conn.Write([]byte(fmt.Sprintf("- %s\n", name)))
	}
}

// broadcast sends a message to all connected clients in the chatroom except the sender.
func broadcast(chatroom *Chatroom, message string, sender net.Conn) {
	chatroom.Mutex.Lock()
	defer chatroom.Mutex.Unlock()

	fmt.Printf("[%s] %s", chatroom.Name, message)
	for client := range chatroom.Clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}

// saveToHistory writes a message to the chatroom's history file.
func saveToHistory(chatroom *Chatroom, message string) {
	if chatroom.History != nil {
		_, err := chatroom.History.WriteString(message)
		if err != nil {
			fmt.Printf("Error writing to history file for chatroom %s: %v\n", chatroom.Name, err)
		}
	}
}
