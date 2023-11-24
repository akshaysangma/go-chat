package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

const (
	Port          = "8125"
	MsgBufferSize = 512
	DebugMode     = false
	RateLimit     = time.Second
	StrikeCount   = 10
	BanTimeout    = 10 * time.Minute
)

func getSensitive(s string) string {
	if !DebugMode {
		return "[REDACTED]"
	}
	return s
}

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	ClientDisconnected
	NewMessage
)

type Message struct {
	msgType MessageType
	conn    net.Conn
	text    string
}

type Client struct {
	conn              net.Conn
	lastMessageSentAt time.Time
	strikeCount       int
}

func server(messages chan Message) {
	clients := make(map[string]*Client)
	bannedClients := make(map[string]time.Time)

	for msg := range messages {
		switch msg.msgType {
		case ClientConnected:
			handleClientConnected(msg, clients, bannedClients)
		case ClientDisconnected:
			handleClientDisconnected(msg, clients)
		case NewMessage:
			handleNewMessage(msg, clients, bannedClients)
		}
	}
}

func connect(conn net.Conn, messages chan Message) {
	messages <- Message{
		conn:    conn,
		msgType: ClientConnected,
	}
	buffer := make([]byte, MsgBufferSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf(
				"Error while reading from client %s, err : %s\n",
				getSensitive(conn.RemoteAddr().String()),
				getSensitive(err.Error()),
			)
			messages <- Message{
				conn:    conn,
				msgType: ClientDisconnected,
			}
			return
		}
		messages <- Message{
			conn:    conn,
			msgType: NewMessage,
			text:    string(buffer[:n]),
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Unable to start TCP listener at port %s, err : %s\n", Port, err)
	}

	log.Printf("Server successfully listening at %s", ln.Addr())
	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept conn, err %s\n", err)
			continue
		}
		go connect(conn, messages)
	}
}

func handleClientConnected(
	msg Message,
	clients map[string]*Client,
	bannedClients map[string]time.Time,
) {
	tcpIP := msg.conn.RemoteAddr().(*net.TCPAddr).IP.String()
	now := time.Now()
	bannedAt, isBanned := bannedClients[tcpIP]
	if isBanned && now.Sub(bannedAt) > BanTimeout {
		delete(bannedClients, tcpIP)
		clients[msg.conn.RemoteAddr().String()].strikeCount = 0
		isBanned = false
	}

	if !isBanned {
		log.Printf("Client %s connected\n", getSensitive(msg.conn.RemoteAddr().String()))
		welcomeMsg := fmt.Sprintf(
			"%s has joined the chat\n",
			getSensitive(msg.conn.RemoteAddr().String()),
		)
		msg.conn.Write([]byte(welcomeMsg))
		clients[msg.conn.RemoteAddr().String()] = &Client{
			conn:              msg.conn,
			lastMessageSentAt: time.Now(),
		}
	} else {
		msg.conn.Write([]byte("You are currently banned!"))
		msg.conn.Close()
	}
}

func handleClientDisconnected(msg Message, clients map[string]*Client) {
	log.Printf("Client %s disconnected\n", getSensitive(msg.conn.RemoteAddr().String()))
	delete(clients, msg.conn.RemoteAddr().String())
	msg.conn.Close()
}

func handleNewMessage(msg Message, clients map[string]*Client, bannedClients map[string]time.Time) {
	addr := msg.conn.RemoteAddr().String()
	now := time.Now()
	client, exists := clients[addr]
	if !exists {
		msg.conn.Write([]byte("You are banned"))
		msg.conn.Close()
		return
	}

	if now.Sub(client.lastMessageSentAt) > RateLimit {
		if utf8.ValidString(msg.text) {
			client.lastMessageSentAt = now
			client.strikeCount = 0
			log.Printf("Client %s sent message: %s", getSensitive(addr), msg.text)
			broadcastMessage(addr, msg.text, clients)
		} else {
			client.strikeCount++
		}
	} else {
		client.strikeCount++
	}

	if client.strikeCount > StrikeCount {
		handleBan(client, addr, bannedClients, clients)
	}
}

func broadcastMessage(senderAddr, text string, clients map[string]*Client) {
	for _, client := range clients {
		if client.conn.RemoteAddr().String() == senderAddr {
			continue
		}
		_, err := client.conn.Write([]byte(text))
		if err != nil {
			log.Printf(
				"Unable to write to %s, err: %s\n",
				getSensitive(client.conn.RemoteAddr().String()),
				getSensitive(err.Error()),
			)
			client.conn.Close()
			delete(clients, client.conn.RemoteAddr().String())
		}
	}
}

func handleBan(
	client *Client,
	addr string,
	bannedClients map[string]time.Time,
	clients map[string]*Client,
) {
	if tcpAddr, ok := client.conn.RemoteAddr().(*net.TCPAddr); ok {
		log.Printf("Client %s banned for excessive messaging", getSensitive(tcpAddr.IP.String()))
		bannedClients[tcpAddr.IP.String()] = time.Now()
		client.conn.Write([]byte("You are banned"))
		client.conn.Close()
		delete(clients, addr)
	}
}
