package main

import (
	"fmt"

	"github.com/swarajroy/toll_calculator/types"
)

type Storer interface {
	Insert(types.Distance) error
	GetDistanceSum(int) (float64, error)
}

type InMemoryStore struct {
	data map[int]float64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[int]float64),
	}
}

func (m *InMemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *InMemoryStore) GetDistanceSum(obuID int) (float64, error) {
	dist, ok := m.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obuID (%d)", obuID)
	}
	return dist, nil
}
