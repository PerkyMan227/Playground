package main

import (
	"fmt"
	"time"
)

type Package struct {
	SYN_flag bool
	ACK_flag bool
	Seq      int
	Ack      int
}

func Client(ISN int, clientToServer chan Package, serverToClient chan Package) {
	SYN := Package{true, false, ISN, 0}
	clientToServer <- SYN
	fmt.Println("Client: SYN sent... Waiting for ACK")
	time.Sleep(time.Second * 2)

	serverPackage := <-serverToClient
	fmt.Println("Client: Recieved ACK+SYN (", serverPackage.Ack, " + ", serverPackage.Seq, ") from server.")
	fmt.Println("Client: Sending ACK to server")
	time.Sleep(time.Second * 2)

	clientToServer <- Package{false, true, serverPackage.Ack, serverPackage.Seq + 1}
}

func Server(ISN int, clientToServer chan Package, serverToClient chan Package) {
	fmt.Println("Server: Server is listening...")
	clientPackge := <-clientToServer

	fmt.Println("Server: Server recieved package! Replying with ACK")
	time.Sleep(time.Second * 2)

	serverToClient <- Package{true, true, ISN, clientPackge.Seq + 1}

	fmt.Println("Server: Listening for final ACK from Client")
	finalPackage := <-clientToServer

	fmt.Println("Server: Server got ::::  ACK_flag:  ", finalPackage.ACK_flag, "SYN_flag: ", finalPackage.SYN_flag, "ack: ", finalPackage.Ack, "Seq: ", finalPackage.Seq)
	fmt.Println("Server: Conncetion established!")
	time.Sleep(time.Second * 2)
}

func main() {
	fmt.Println("=======================")

	clientToServer := make(chan Package)
	serverToClient := make(chan Package)

	go Server(5000, clientToServer, serverToClient)
	go Client(1000, clientToServer, serverToClient)

	select {}
}
