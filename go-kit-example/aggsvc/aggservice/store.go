package aggservice

import (
	"fmt"

	"github.com/yuriykis/tolling/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
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

func (ms *MemoryStore) Get(obuID int) (float64, error) {
	dist, ok := ms.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for OBU %d", obuID)
	}
	return dist, nil
}
