package main

import (
	"fmt"

	"github.com/abdoroot/tolling/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type mermoryStore struct {
	data map[int]float64
}

func NewMemoryStore() Storer {
	return &mermoryStore{
		data: make(map[int]float64, 0),
	}
}

func (s *mermoryStore) Insert(data types.Distance) error {
	s.data[data.OBUID] += data.Value
	return nil
}

func (s *mermoryStore) Get(OBUID int) (float64, error) {
	dist, ok := s.data[OBUID]
	if !ok {
		return 0, fmt.Errorf("distance not found for obuid %v", OBUID)
	}
	return dist, nil
}
