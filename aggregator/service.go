package main

import "github.com/abdoroot/tolling/types"

const priceRate = 2.26

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
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

func (i *invoiceAggregator) CalculateInvoice(obuid int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuid)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuid,
		TotalDistance: dist,
		TotalAmount:   dist * priceRate,
	}
	return inv, nil
}
