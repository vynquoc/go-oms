package main

import (
	"errors"
	"net/http"

	common "github.com/vynquoc/go-oms-common"
	pb "github.com/vynquoc/go-oms-common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	errorStatus := status.Convert(err)
	if errorStatus != nil {
		if errorStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, errorStatus.Message())
			return
		}
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJson(w, http.StatusOK, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItem
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("Item ID is required")
		}

		if item.Quantity <= 0 {
			return errors.New("Item must have a valid quantity")
		}
	}

	return nil
}
