package gateway

import (
	"context"
	"log"

	pb "github.com/vynquoc/go-oms-common/api"
	"github.com/vynquoc/go-oms-common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrdersServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: payload.CustomerID,
		Items:      payload.Items,
	})
}

func (g *gateway) GetOrder(ctx context.Context, orderID, customerID string) (*pb.Order, error) {

	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrdersServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderID:    orderID,
		CustomerID: customerID,
	})
}
