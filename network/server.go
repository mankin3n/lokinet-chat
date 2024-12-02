package network

import (
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

	// Get the public IP address
	publicIP := getPublicIP()
	if publicIP == "" {
		publicIP = "127.0.0.1" // Fallback to localhost
	}

	// Display connection details
	fmt.Printf("Chatroom server started!\n")
	fmt.Printf("Clients can connect using the following details:\n")
	fmt.Printf("Public IP Address: %s\n", publicIP)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Example: %s:%s\n\n", publicIP, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleClient(conn)
	}
}

// getPublicIP retrieves the public IP address of the server.
func getPublicIP() string {
	// Retrieve public IP from a web service
	resp, err := net.LookupIP("myip.opendns.com")
	if err != nil {
		return ""
	}

	for _, ip := range resp {
		if strings.Contains(ip.String(), ".") { // IPv4 only
			return ip.String()
		}
	}
	return ""
}

// handleClient handles individual client connections.
func handleClient(conn net.Conn) {
	defer conn.Close()
	// Handle client logic (joining chatrooms, sending/receiving messages)
}
