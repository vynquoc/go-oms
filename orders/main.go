package main

import (
	"context"
	"log"
	"net"

	common "github.com/vynquoc/go-oms-common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:3001")
)

func main() {
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen %v", grpcAddr)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)

	NewGRPCHandler(grpcServer)
	svc.CreateOrder(context.Background())

	log.Println("GRPC Order service started at:", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
