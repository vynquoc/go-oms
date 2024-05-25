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

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID:         "69",
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}

	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, payload *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(payload.Items) == 0 {
		return nil, common.ErrNoItem
	}

	mergedItems := mergeItemsQuantities(payload.Items)
	fmt.Println(mergedItems)

	var itemsWithPrice []*pb.Item

	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			ID:       i.ID,
			PriceID:  "price_1PKLtn2KofBMjvPrmKWcckJS",
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
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
