package inmem

import pb "github.com/vynquoc/go-oms-common/api"

type InMem struct {
}

func NewInmem() *InMem {
	return &InMem{}
}

func (i *InMem) CreatePaymentLink(*pb.Order) (string, error) {
	return "dummylink", nil
}
