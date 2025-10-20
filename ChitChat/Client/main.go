package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.NewClient("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()

	client := proto.NewChitChatServiceClient(connection)

	response, err := client.Join(context.Background(), &proto.JoinRequest{
		ParticipantName: "Alice",
	})
	if err != nil {
		log.Fatalf("Failed to join: %v", err)
	}

	log.Printf("Joined with ID: %s", response.ParticipantId)
	log.Printf("Response from server: %s", response.Message)

	myID := response.ParticipantId
	_ = myID

	//hurtig test her, Server modtager det fint.
	client.Publish(context.Background(), &proto.PublishRequest{ParticipantId: myID, Content: "Hello from Allice at home! I hope everyone receives this message)"})

}

func receiveBroadcasts(client proto.ChitChatServiceClient) {
	stream, err := client.ReceiveBroadcasts(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Failed to receive broadcasts: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server")
			break
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		fmt.Printf("\n[%d] %s: %s\n> ", msg.LamportTimestamp, msg.Sender, msg.Content)
	}
}
