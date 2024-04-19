package main

import (
	"fmt"

	"github.com/swarajroy/toll_calculator/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
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
