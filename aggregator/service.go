package main

import "github.com/abdoroot/tolling/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type invoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &invoiceAggregator{
		store: store,
	}
}

func (i *invoiceAggregator) AggregateDistance(data types.Distance) error {
	return i.store.Insert(data)
}
