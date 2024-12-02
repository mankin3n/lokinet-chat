
# Lokinet Chat CLI

Lokinet Chat CLI is a secure, multi-chatroom messaging application built with Go. It allows users to create profiles, join or create chatrooms, and exchange messages in real-time, with chat histories saved on the server.

---

## Features

1. **Multi-Chatroom Support**
   - Users can create or join multiple chatrooms hosted on a server.

2. **Real-Time Messaging**
   - Messages are broadcast to all participants in the chatroom.

3. **User Profiles**
   - Users can create and manage profiles with public/private keys.

4. **Chat History**
   - All chat messages are saved in history files on the server.

---

## Prerequisites

- Go installed on your system (version 1.20 or higher).
- Access to the Lokinet Chat project files.

---

## Installation

1. **Clone the Repository**
   ```bash
   git clone <repository_url>
   cd lokinet-chat
   ```

2. **Install Dependencies**
   Ensure you have Go modules initialized:
   ```bash
   go mod tidy
   ```

3. **Build the Application**
   Compile the application:
   ```bash
   go build -o lokinet-chat .
   ```

---

## Usage

1. **Start the Chatroom Server**
   - Run the application and choose `Option 1` to start the server:
     ```bash
     ./lokinet-chat
     ```
   - The server will listen on `localhost:8080` by default.

2. **Create or Join a Chatroom**
   - Run the application on another terminal or machine and choose `Option 2`.
   - Enter the server address and the chatroom name to join or create a chatroom.

3. **Exchange Messages**
   - Messages sent by one user in a chatroom will be broadcast to all other users.

4. **Manage Profiles**
   - Use `Option 3` to view your profile.
   - Use `Option 4` to create a new profile.

---

## Chat History

- Chat histories are saved in files named `<chatroom_name>_history.txt` on the server.
- Each chatroom maintains its own separate history.

---

## Exit the Application

- To leave a chatroom, type `/quit`.
- To exit the application, select `Option 5` from the main menu.

---

## Future Enhancements

- Private messaging between users.
- Chatroom access permissions.
- Persistent user data storage.

---

## License

This project is licensed under the MIT License.

---

## Author

Created by **Lokinet Chat CLI Team**.

