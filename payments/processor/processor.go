package processor

import pb "github.com/vynquoc/go-oms-common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
