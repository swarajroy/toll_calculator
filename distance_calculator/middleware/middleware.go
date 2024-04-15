package middleware

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/distance_calculator/service"
	"github.com/swarajroy/toll_calculator/types"
)

type LogMiddleware struct {
	next service.CalculatorServicer
}

func NewLogMiddleware(next service.CalculatorServicer) service.CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data *types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info("calculate distance")
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return
}
