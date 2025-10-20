package main

import (
	proto "ChitChat/grpc"
	"context"
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
	log.Printf("MESSGE: %s", response.Message)

}
