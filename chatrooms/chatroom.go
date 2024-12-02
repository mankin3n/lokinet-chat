package chatrooms

import "fmt"

func CreateChatRoom(name string) {
	fmt.Println("Chatroom created:", name)
}

func JoinChatRoom(name string) {
	fmt.Println("Joined chatroom:", name)
}

func LeaveChatRoom(name string) {
	fmt.Println("Left chatroom:", name)
}
