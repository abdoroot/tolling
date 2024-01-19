package main

import (
	"time"

	"github.com/abdoroot/tolling/types"
	"github.com/sirupsen/logrus"
)

type logMiddleWare struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &logMiddleWare{
		next: next,
	}
}

func (l *logMiddleWare) AggregateDistance(data types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuid": data.OBUID,
			"took":  time.Since(start),
			"err":   err,
		}).Info("aggregate")
	}(time.Now())
	err = l.next.AggregateDistance(data)
	return err
}

func (l *logMiddleWare) CalculateInvoice(OBUID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("CalculateInvoice")
	}(time.Now())
	inv, err = l.next.CalculateInvoice(OBUID)
	return
}
