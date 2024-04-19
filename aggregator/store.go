package main

import "github.com/swarajroy/toll_calculator/types"

type Storer interface {
	Insert(types.Distance) error
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
