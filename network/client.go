package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// JoinChatroom connects to the specified chatroom on the server.
func JoinChatroom(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn)

	fmt.Println("Connected to the server. Follow the instructions to join or create a chatroom.")

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
