package main

import (
	"time"

	"github.com/abdoroot/tolling/types"
	"github.com/sirupsen/logrus"
)

type logMiddleware struct {
	next CalculateServicer
}

func NewLogMiddleware(next CalculateServicer) CalculateServicer {
	return &logMiddleware{
		next: next,
	}
}

func (m *logMiddleware) CalculateDistance(data types.OBUdata) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took": time.Since(start),
				"err":  err,
				"dist": dist,
			},
		).Info("Calculate Distance")
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return dist, err
}
