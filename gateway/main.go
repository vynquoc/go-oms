package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vynquoc/go-oms-common"
	pb "github.com/vynquoc/go-oms-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr             = common.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = "localhost:3001"
)

func main() {
	conn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to dial grpc at %v", orderServiceAddr)
	}
	defer conn.Close()
	c := pb.NewOrdersServiceClient(conn)

	log.Println("Dialing order service at", orderServiceAddr)
	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Failed to create server")
	}
}
