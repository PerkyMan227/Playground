package main

import (
	proto "ChitChat/grpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

// ChitChatService implements the gRPC server interface.
type Server struct {
	grpcServer *grpc.Server
	id         string
	proto.UnimplementedChitChatServiceServer
}

func main() {
	server := &Server{}
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
		log.Fatalf("Failed to serve: %v", err)
	}
}
