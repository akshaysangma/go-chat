# TCP Chat Server with Go Std Lib

## Description
This project is a simple TCP chat server implemented in Go using standard library. It features handling of client connections, message broadcasting, rate limiting, and client banning based on strike counts. The server operates over TCP and provides basic chat functionality.

## Features
- TCP-based communication
- Handling multiple client connections
- Broadcasting messages to all connected clients
- Rate limiting to prevent spam
- Strike-based banning mechanism
- Configurable settings

## Getting Started

### Prerequisites
- Go (written and tested in 1.21.x)

### Installation
1. Clone the repository:

```shell
git clone https://github.com/akshaysangma/go-chat.git 
```
2. Navigate to the project directory:
```shell
cd go-chat
```

### Configuration
The server configuration can be adjusted via the following constants in the main Go file:
- `Port`: The port on which the server listens.
- `MsgBufferSize`: Size of the message buffer.
- `DebugMode`: Enables or disables debug mode.
- `RateLimit`: Duration for rate limiting messages (in seconds).
- `StrikeCount`: Number of strikes before banning a client.
- `BanTimeout`: Duration for which a client is banned (in minutes).

### Running the Server
To start the server, run:
```shell
go run main.go
```

## Usage
Once the server is running, clients can connect to it using any TCP client at the specified port. Messages sent by a client will be broadcast to all other connected clients.
Example with Telnet
```shell
telnet localhost 8125
```

## Contributing
Contributions to the project are welcome. Please follow the standard Git workflow - fork the repository, make your changes, and submit a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

