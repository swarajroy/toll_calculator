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

func (m *LogMiddleware) Invoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if invoice != nil {
			distance = invoice.TotalDistance
			amount = invoice.InvoiceAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"distance": distance,
			"amount":   amount,
		}).Info("Invoice")
	}(time.Now())
	invoice, err = m.next.Invoice(obuID)
	return
}
