package main

import (
	"fmt"
	"sync"
)

type Message struct {
	text string
	SYN  int
	ACK  int
}

func client(server, client chan Message, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("client")

	clientSYN := 100
	clientACK := 0

	server <- Message{SYN: clientSYN}

	for msg := range client {
		if msg.text == "" {
			// Dette er et handshake
			clientACK = msg.SYN + 1
			// Sender sidste handshake
			fmt.Println("Client received server SYN:", msg.SYN, "setting client ACK to:", clientACK)
			fmt.Println("Client sending final SYN+ACK:", clientSYN+1, "ACK:", clientACK)
			server <- Message{SYN: clientSYN + 1, ACK: clientACK}
		} else {
			fmt.Println("Message from server: ", msg.text)
		}
	}
}

func server(client, server chan Message, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("server")
	serverSYN := 200
	serverACK := 0

	for msg := range server {
		if msg.text == "" && msg.ACK == 0 {
			// Dette er et handshake
			serverACK = msg.SYN + 1
			fmt.Println("Server received client SYN:", msg.SYN, "setting server ACK to:", serverACK)
			fmt.Println("Server sending SYN:", serverSYN, "ACK:", serverACK)
			client <- Message{SYN: serverSYN, ACK: serverACK}
		}
		if msg.text == "" && msg.ACK == serverSYN+1 {
			fmt.Println("Server received client SYN+ACK:", msg.SYN, "ACK:", msg.ACK)
			fmt.Println("Server sending message with SYN:", serverSYN+1, "ACK:", serverACK+1)
			client <- Message{SYN: serverSYN + 1, ACK: serverACK + 1, text: "What do you want??"}
			close(server)
			close(client)
		}
	}
}

func main() {
	fmt.Println("Hello Worldies!")
	clientChan := make(chan Message)
	serverChan := make(chan Message)

	var wg sync.WaitGroup
	wg.Add(2)

	go client(serverChan, clientChan, &wg)
	go server(clientChan, serverChan, &wg)

	wg.Wait()

}
