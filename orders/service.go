package main

import (
	"context"
	"fmt"

	common "github.com/vynquoc/go-oms-common"
	pb "github.com/vynquoc/go-oms-common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, payload *pb.CreateOrderRequest) error {
	if len(payload.Items) == 0 {
		return common.ErrNoItem
	}

	mergedItems := mergeItemsQuantities(payload.Items)
	fmt.Println(mergedItems)

	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false

		for _, mergedItem := range merged {
			if mergedItem.ID == item.ID {
				mergedItem.Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
