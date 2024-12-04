package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// JoinChatroom connects to the specified chatroom server and handles interaction.
func JoinChatroom(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	// Display server messages
	go receiveMessages(conn)

	// Read user input and send to server
	fmt.Println("Connected to the server. Follow the instructions to join or create a chatroom.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "/quit" {
			fmt.Println("Exiting chatroom...")
			break
		}
		fmt.Fprintf(conn, "%s\n", message)
	}
}

// receiveMessages reads messages from the server and displays them.
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server.")
			return
		}
		fmt.Print(message)
	}
}
