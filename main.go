package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"lokinet-chat/network"
	"lokinet-chat/user"
)

// ANSI Escape Codes for clearing screen
const clearScreen = "\033[H\033[2J"

func printMenu() {
	fmt.Print(clearScreen)
	fmt.Println(`
 _  _ _  _ _    ___  ____ ____    _  _ ____ ____ ____ ____ ____ ____ ____ ____ 
|__| |  | |    |  \ |___ |__/    |\/| |___ [__  [__  |__| | __ |___ |__/ 
|  | |__| |___ |__/ |___ |  \    |  | |___ ___] ___] |  | |__] |___ |  \ 
`)
	fmt.Println("[1] Start Chatroom Server")
	fmt.Println("[2] Join or Create a Chatroom")
	fmt.Println("[3] View Profile")
	fmt.Println("[4] Create Profile")
	fmt.Println("[5] Exit")
	fmt.Println("---------------------------------------------------")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
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
			network.JoinChatroom(address)
		case "3":
			// Display Profile Without Clearing Screen
			fmt.Println("=== View Profile ===")
			profile, err := user.LoadProfile()
			if err != nil {
				fmt.Println("Error: Profile not found! Please create one first.")
				fmt.Println("Press Enter to return to the menu.")
				reader.ReadString('\n') // Wait for user input before returning
				continue
			}

			fmt.Printf("Username: %s\n", profile.Username)
			fmt.Printf("Public Key:\n%s\n", profile.PublicKey)
			if len(profile.PrivateKey) > 20 {
				fmt.Printf("Private Key (masked):\n%s...%s\n", profile.PrivateKey[:10], profile.PrivateKey[len(profile.PrivateKey)-10:])
			} else {
				fmt.Printf("Private Key:\n%s\n", profile.PrivateKey)
			}
			fmt.Println("Press Enter to return to the menu.")
			reader.ReadString('\n') // Wait for user input before returning
		case "4":
			fmt.Print("Enter your username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			existingProfile, err := user.LoadProfile()
			if err == nil && existingProfile.Username == username {
				fmt.Printf("Error: A profile with username '%s' already exists.\n", username)
				fmt.Println("Press Enter to return to the menu.")
				reader.ReadString('\n') // Wait for user input before returning
				continue
			}

			profile, err := user.CreateProfile(username)
			if err != nil {
				fmt.Println("Error creating profile:", err)
				fmt.Println("Press Enter to return to the menu.")
				reader.ReadString('\n') // Wait for user input before returning
				continue
			}

			err = user.SaveProfile(profile)
			if err != nil {
				fmt.Println("Error saving profile:", err)
			} else {
				fmt.Println("Profile created successfully!")
			}
			fmt.Println("Press Enter to return to the menu.")
			reader.ReadString('\n') // Wait for user input before returning
		case "5":
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
