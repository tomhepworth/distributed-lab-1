package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	errormsg := err.Error()
	fmt.Printf(errormsg)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	for {
		conn, err := ln.Accept()
		fmt.Println("Hello new client")
		if err != nil {
			handleError(err)
		}
		// and add it to the channel for handling connections.
		conns <- conn
	}

}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	fmt.Println("handling", clientid)

	reader := bufio.NewReader(client)
	for {
		///fmt.Println("bbbb")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Died:", clientid)
			handleError(err)
			break
		}
		tidy := msg[:len(msg)-1]

		finalMsg := Message{sender: clientid, message: tidy}
		msgs <- finalMsg
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", "127.0.0.1:8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)

	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			cn := len(clients)
			clients[cn] = conn
			// - start to asynchronously handle messages from this client
			go handleClient(clients[cn], cn, msgs)

		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			fmt.Println("aaaaa")
			fmt.Println(clients)
			for i, conn := range clients {
				fmt.Println(i)
				if i != msg.sender {
					fmt.Println("Sending: " + msg.message)
					fmt.Fprintln(conn, msg.message)
				}
			}
		}
	}
}
