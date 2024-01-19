package main

import "github.com/abdoroot/tolling/types"

type Storer interface {
	Insert(types.Distance) error
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
