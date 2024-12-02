package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// JoinChatroom connects to the specified chatroom on the server.
func JoinChatroom(address string, chatroomName string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	// Send chatroom name
	fmt.Fprintf(conn, "%s\n", chatroomName)

	go receiveMessages(conn)

	fmt.Println("Connected to chatroom:", chatroomName)
	fmt.Println("Start typing your messages. Type '/quit' to exit.")

	// Read user input and send to the server
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
			fmt.Println("Disconnected from chatroom.")
			return
		}
		fmt.Print(message)
	}
}
