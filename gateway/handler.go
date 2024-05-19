package main

import (
	"net/http"

	common "github.com/vynquoc/go-oms-common"
	pb "github.com/vynquoc/go-oms-common/api"
)

type handler struct {
	client pb.OrdersServiceClient
}

func NewHandler(client pb.OrdersServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity

	if err := common.ReadJson(r, &items); err != nil {
		common.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
}
