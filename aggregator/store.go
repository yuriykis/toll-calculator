package main

import "github.com/yuriykis/tolling/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() Storer {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (ms *MemoryStore) Insert(d types.Distance) error {
	ms.data[d.OBUID] += d.Value
	return nil
}
