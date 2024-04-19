package main

import (
	"fmt"
	"math/rand"

	"github.com/swarajroy/toll_calculator/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
	Invoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in storage: ", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) Invoice(obuID int) (*types.Invoice, error) {
	distSum, err := i.store.GetDistanceSum(obuID)
	if err != nil {
		return nil, err
	}
	return types.NewInvoice(obuID, distSum, rand.Float64()*distSum), nil
}
