package grpc

import (
	"context"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"chat-service/internal/grpc/chatpb"
	"chat-service/internal/service"

	"google.golang.org/grpc"
)

// ChatServer - gRPC сервер
type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
	chatService *service.ChatService
}

// NewChatServer - конструктор
func NewChatServer(chatService *service.ChatService) *ChatServer {
	return &ChatServer{chatService: chatService}
}

// CreateChat - обработка gRPC-запроса
func (s *ChatServer) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.CreateChatResponse, error) {
	return s.chatService.CreateChat(ctx, req)
}

func RunGRPCServer(chatService *service.ChatService) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Ошибка запуска gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, NewChatServer(chatService))

	reflection.Register(grpcServer)

	log.Println("🚀 Chat Service (gRPC) запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Ошибка gRPC-сервера: %v", err)
	}
}
