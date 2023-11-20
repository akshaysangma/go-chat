package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

const (
	PORT          = "8125"
	MSGBUFFERSIZE = 512
	DEBUGMODE     = false
	RATELIMIT     = 1.0
	STRIKECOUNT   = 10
	BANTIMEOUT    = 10.0 * 60.0
)

func getSensitive(s string) string {
	if !DEBUGMODE {
		return "[REDACTED]"
	}
	return s
}

type MessageType int

const (
	clientConnected MessageType = iota + 1
	clientDisconnected
	newMessage
)

type Message struct {
	msgType MessageType
	conn    net.Conn
	text    string
}

type Client struct {
	conn              net.Conn
	lastMessageSentAt time.Time
	strikecount       int
}

func server(messages chan Message) {
	clients := make(map[string]*Client)
	bannedClients := make(map[string]time.Time)
	for {
		msg := <-messages
		switch msg.msgType {
		case clientConnected:
			tcpIP := msg.conn.RemoteAddr().(*net.TCPAddr).IP.String()
			now := time.Now()
			bannedAt, isbanned := bannedClients[tcpIP]
			if isbanned {
				if now.Sub(bannedAt).Seconds() > BANTIMEOUT {
					delete(bannedClients, tcpIP)
					clients[msg.conn.RemoteAddr().String()].strikecount = 0
					isbanned = false
				}
			}
			if !isbanned {
				log.Printf("Client %s connected\n", getSensitive(msg.conn.RemoteAddr().String()))
				welcomemsg := fmt.Sprintf(
					"%s has join the chat\n",
					getSensitive(msg.conn.RemoteAddr().String()),
				)
				msg.conn.Write([]byte(welcomemsg))
				clients[msg.conn.RemoteAddr().String()] = &Client{
					conn:              msg.conn,
					lastMessageSentAt: time.Now(),
				}
			} else {
				msg.conn.Write([]byte("You are currently banned!"))
				msg.conn.Close()
			}
		case clientDisconnected:
			log.Printf("Client %s disconnected\n", getSensitive(msg.conn.RemoteAddr().String()))
			delete(clients, msg.conn.RemoteAddr().String())
			msg.conn.Close()
		case newMessage:
			addr := msg.conn.RemoteAddr().String()
			now := time.Now()
			if clients[addr] != nil {
				if now.Sub(clients[addr].lastMessageSentAt).Seconds() > RATELIMIT {
					if utf8.ValidString(msg.text) {
						clients[addr].lastMessageSentAt = now
						clients[addr].strikecount = 0
						log.Printf(
							"Client %s sent message %s",
							getSensitive(msg.conn.RemoteAddr().String()),
							msg.text,
						)
						for _, client := range clients {
							if client.conn.RemoteAddr().String() == addr {
								continue
							}
							_, err := client.conn.Write([]byte(msg.text))
							if err != nil {
								log.Printf(
									"Unable to write to %s, err : %s/n",
									getSensitive(client.conn.RemoteAddr().String()),
									getSensitive(err.Error()),
								)
								messages <- Message{conn: client.conn, msgType: clientDisconnected}
							}
						}
					} else {
						clients[addr].strikecount++
					}
				} else {
					clients[addr].strikecount++
				}
				if clients[addr].strikecount > STRIKECOUNT {
					if tcpAddr, ok := clients[addr].conn.RemoteAddr().(*net.TCPAddr); ok {
						log.Printf(
							"Client %s banned for too much message",
							getSensitive(tcpAddr.IP.String()),
						)
						bannedClients[tcpAddr.IP.String()] = now
						msg.conn.Write([]byte("You are banned"))
						clients[addr].conn.Close()
						delete(clients, addr)
					}
				}
			} else {
				msg.conn.Write([]byte("You are banned"))
				msg.conn.Close()
			}
		}
	}
}

func connect(conn net.Conn, messages chan Message) {
	messages <- Message{
		conn:    conn,
		msgType: clientConnected,
	}
	buffer := make([]byte, MSGBUFFERSIZE)
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
				msgType: clientDisconnected,
			}
			return
		}
		messages <- Message{
			conn:    conn,
			msgType: newMessage,
			text:    string(buffer[:n]),
		}
	}
}

func main() {
	// TODO : Environemnt Variables
	// TODO : replace log with slog

	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("Unable to start TCP listener at port %s, err : %s\n", PORT, err)
	}

	log.Printf("Server successfully listening at %s", ln.Addr())
	messages := make(chan Message)
	go server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept conn, err %s\n", err)
		}
		go connect(conn, messages)
	}
}
