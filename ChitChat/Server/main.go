package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// ChitChatService implements the gRPC server interface.
type Server struct {
	grpcServer *grpc.Server
	id         string
	mu         sync.Mutex
	proto.UnimplementedChitChatServiceServer

	streams          map[string]proto.ChitChatService_ReceiveBroadcastsServer
	participants     map[string]string
	lamportTimestamp int64
	nextID           int
}

func main() {
	server := &Server{
		participants: make(map[string]string),
		nextID:       1,
	}
	server.StartServer()
}

func (s *Server) StartServer() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	proto.RegisterChitChatServiceServer(s.grpcServer, s)

	if err := s.grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}

func (s *Server) Join(ctx context.Context, req *proto.JoinRequest) (*proto.JoinResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	participantID := fmt.Sprintf("user_%d", s.nextID)
	s.nextID++

	s.participants[participantID] = req.ParticipantName
	log.Printf("%s joined (ID: %s)", req.ParticipantName, participantID)

	s.lamportTimestamp++
	s.broadcastLocked(&proto.BroadcastMessage{
		Content:          fmt.Sprintf("%s just joined the chat! Say hi!", req.ParticipantName),
		LamportTimestamp: s.lamportTimestamp,
		Sender:           "le system",
	})

	return &proto.JoinResponse{
		ParticipantId: participantID,
		Message:       "Request to join accepted! Welcome to Freaky Chat!",
	}, nil

}

func (s *Server) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.PublishResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name, exists := s.participants[req.ParticipantId]
	if !exists {
		return &proto.PublishResponse{
				Success: false,
			},
			fmt.Errorf("Participant not found. Request unsuccessfull...")
	}
	s.lamportTimestamp++

	msg := &proto.BroadcastMessage{
		Content:          req.Content,
		LamportTimestamp: s.lamportTimestamp,
		Sender:           name,
	}
	log.Printf("[%d] %s: %s", msg.LamportTimestamp, msg.Sender, msg.Content)

	s.broadcastLocked(msg)

	return &proto.PublishResponse{Success: true}, nil
}

func (s *Server) broadcastLocked(msg *proto.BroadcastMessage) {
	for streamID, stream := range s.streams {
		if err := stream.Send(msg); err != nil {
			log.Printf("Failed to send to %s: %v", streamID, err)
			delete(s.streams, streamID)
		}
	}
}
