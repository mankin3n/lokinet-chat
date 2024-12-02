package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"lokinet-chat/network"
	"lokinet-chat/user"
)

func main() {
	fmt.Println("=== Lokinet Chat CLI ===")
	fmt.Println("[1] Start Chatroom Server")
	fmt.Println("[2] Join a Chatroom")
	fmt.Println("[3] View Profile")
	fmt.Println("[4] Create Profile")
	fmt.Println("[5] Exit")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Select an option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Println("Starting the chatroom server...")
			go network.StartServer()
		case "2":
			fmt.Print("Enter the server address (e.g., 127.0.0.1:8080): ")
			address, _ := reader.ReadString('\n')
			address = strings.TrimSpace(address)

			fmt.Print("Enter the chatroom name: ")
			chatroomName, _ := reader.ReadString('\n')
			chatroomName = strings.TrimSpace(chatroomName)

			network.JoinChatroom(address, chatroomName)
		case "3":
			fmt.Println("=== View Profile ===")
			profile, err := user.LoadProfile()
			if err != nil {
				fmt.Println("Error: Profile not found! Please create one first.")
				continue
			}

			fmt.Printf("Username: %s\n", profile.Username)
			fmt.Printf("Public Key:\n%s\n", profile.PublicKey)
			if len(profile.PrivateKey) > 20 {
				fmt.Printf("Private Key (masked):\n%s...%s\n", profile.PrivateKey[:10], profile.PrivateKey[len(profile.PrivateKey)-10:])
			} else {
				fmt.Printf("Private Key:\n%s\n", profile.PrivateKey)
			}
		case "4":
			fmt.Print("Enter your username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			existingProfile, err := user.LoadProfile()
			if err == nil && existingProfile.Username == username {
				fmt.Printf("Error: A profile with username '%s' already exists.\n", username)
				continue
			}

			profile, err := user.CreateProfile(username)
			if err != nil {
				fmt.Println("Error creating profile:", err)
				continue
			}

			err = user.SaveProfile(profile)
			if err != nil {
				fmt.Println("Error saving profile:", err)
			} else {
				fmt.Println("Profile created successfully!")
			}
		case "5":
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
