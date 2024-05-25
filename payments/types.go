package main

import (
	"context"
	pb "github.com/vynquoc/go-oms-common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
