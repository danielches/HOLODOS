package main

import (
	"HOLODOS/internal/server"
	"HOLODOS/internal/service"
	"HOLODOS/internal/storage"
	fridgev1 "HOLODOS/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Инициализация зависимостей
	storage := storage.NewMemoryStorage()
	service := service.NewFridgeService(storage)
	fridgeServer := server.NewFridgeServer(service)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	fridgev1.RegisterFridgeServiceServer(grpcServer, fridgeServer)

	// Запуск сервера
	lis, err := net.Listen("tcp", ":44044")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on :44044")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
