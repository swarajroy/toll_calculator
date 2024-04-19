package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(agg Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: agg,
	}
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(dist)
	return
}
