package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/vynquoc/go-oms-common/api"
	"github.com/vynquoc/go-oms-common/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrdersServiceServer
	service OrdersService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrdersServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received %v", payload)
	items, err := h.service.ValidateOrder(ctx, payload)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, payload, items)
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        marshalledOrder,
	})
	return o, nil
}

func (h *grpcHandler) GetOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, payload)

}
