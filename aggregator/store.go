package main

import "github.com/yuriykis/tolling/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
}

func NewMemoryStore() Storer {
	return &MemoryStore{}
}

func (ms *MemoryStore) Insert(d types.Distance) error {
	return nil
}
