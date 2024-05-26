package main

import (
	"context"
	"testing"

	"github.com/vynquoc/go-oms-common/api"
	"github.com/vynquoc/go-oms-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error: %v", err)
		}

		if link == "" {
			t.Errorf("CreatePayment() link is empty")
		}
	})
}
