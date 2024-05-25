package main

import (
	"context"

	pb "github.com/vynquoc/go-oms-common/api"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(context.Context, *pb.Order) (string, error) {
	//connect to payment processor

	return "", nil
}
