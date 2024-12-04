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
	LogChannel  = make(chan string, 100)     // Exported log channel
)

// Highlight text for username (using ANSI escape codes for bold text)
func highlightUsername(username string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", username)
}

// StartServer initializes the chatroom server.
func StartServer() {
	port := "8080" // Default port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// Get the public IP address using internal tools
	publicIP := getPublicIP()
	if publicIP == "" {
		publicIP = "127.0.0.1" // Fallback to localhost
	}

	// Display connection details
	LogChannel <- fmt.Sprintf("Chatroom server started!\n")
	LogChannel <- fmt.Sprintf("Clients can connect using the following details:\n")
	LogChannel <- fmt.Sprintf("Public IP Address: %s\n", publicIP)
	LogChannel <- fmt.Sprintf("Port: %s\n", port)
	LogChannel <- fmt.Sprintf("Example: %s:%s\n\n", publicIP, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogChannel <- fmt.Sprintf("Connection error: %v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

// getPublicIP retrieves the public IP address of the server using socket connection.
func getPublicIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // Use Google's public DNS server for a connection test
	if err != nil {
		LogChannel <- fmt.Sprintf("Error determining public IP: %v\n", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// handleClient manages communication with a single client.
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Enter username
	conn.Write([]byte("Enter your username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	LogChannel <- fmt.Sprintf("%s connected.\n", highlightUsername(username))
	conn.Write([]byte("Welcome to the server, " + username + "!\n"))

	for {
		// List chatrooms
		chatroomsMu.Lock()
		conn.Write([]byte("Available chatrooms:\n"))
		for name := range chatrooms {
			conn.Write([]byte(fmt.Sprintf("- %s\n", name)))
		}
		conn.Write([]byte("Enter a chatroom name to join or create: "))
		chatroomsMu.Unlock()

		chatroomName, _ := reader.ReadString('\n')
		chatroomName = strings.TrimSpace(chatroomName)

		chatroom := getOrCreateChatroom(chatroomName)

		chatroom.Mutex.Lock()
		chatroom.Clients[conn] = username
		chatroom.Mutex.Unlock()

		LogChannel <- fmt.Sprintf("%s joined chatroom: %s\n", highlightUsername(username), chatroomName)
		broadcast(chatroom, fmt.Sprintf("%s has joined the chatroom\n", username), conn)

		// Chat loop
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				chatroom.Mutex.Lock()
				delete(chatroom.Clients, conn)
				chatroom.Mutex.Unlock()
				broadcast(chatroom, fmt.Sprintf("%s has left the chatroom\n", username), conn)
				LogChannel <- fmt.Sprintf("%s left chatroom: %s\n", highlightUsername(username), chatroomName)
				break
			}

			message = strings.TrimSpace(message)
			if strings.ToLower(message) == "/exit" {
				conn.Write([]byte("Exiting chatroom...\n"))
				chatroom.Mutex.Lock()
				delete(chatroom.Clients, conn)
				chatroom.Mutex.Unlock()
				broadcast(chatroom, fmt.Sprintf("%s has left the chatroom\n", username), conn)
				LogChannel <- fmt.Sprintf("%s exited chatroom: %s\n", highlightUsername(username), chatroomName)
				break
			}

			fullMessage := fmt.Sprintf("%s: %s\n", username, message)
			broadcast(chatroom, fullMessage, conn)
			LogChannel <- fmt.Sprintf("[%s] %s\n", chatroomName, highlightUsername(username)+": "+message)
		}
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
	LogChannel <- fmt.Sprintf("Created new chatroom: %s\n", name)
	return newChatroom
}

// broadcast sends a message to all connected clients in the chatroom except the sender.
func broadcast(chatroom *Chatroom, message string, sender net.Conn) {
	chatroom.Mutex.Lock()
	defer chatroom.Mutex.Unlock()

	for client := range chatroom.Clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}
