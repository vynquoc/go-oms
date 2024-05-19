package main

import (
	"context"
	"log"

	pb "github.com/vynquoc/go-oms-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrdersServiceServer
}

func NewGRPCHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrdersServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received %v", payload)
	o := &pb.Order{
		ID: "9",
	}
	return o, nil
}
