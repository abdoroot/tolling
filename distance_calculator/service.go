package main

import (
	"math"

	"github.com/abdoroot/tolling/types"
)

type CalculateServicer interface {
	CalculateDistance(data types.OBUdata) (float64, error)
}

type calculateService struct {
	PrevPoint []float64
}

func NewCalculateService() CalculateServicer {
	return &calculateService{}
}

func (s *calculateService) CalculateDistance(data types.OBUdata) (float64, error) {
	distance := 0.0
	if len(s.PrevPoint) > 0 {
		distance = calculateDistance(s.PrevPoint[0], s.PrevPoint[1], data.Lat, data.Long)
	}
	s.PrevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
