package gateway

import (
	"context"

	pb "github.com/vynquoc/go-oms-common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
