package main

import (
	"github.com/abdoroot/tolling/types"
	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	Next DataProducer
}

func NewLoggingMiddleware(next DataProducer) *LoggingMiddleware {
	return &LoggingMiddleware{
		Next: next,
	}
}

//implement DataProducer interface
func (m *LoggingMiddleware) ProduceData(data types.OBUdata) error {
	logrus.WithFields(
		logrus.Fields{
			"obu_id": data.OBUID,
			"lat":    data.Lat,
			"long":   data.Long,
		},
	).Info("Produce")
	return m.Next.ProduceData(data)
}
