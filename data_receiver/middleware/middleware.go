package middleware

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/data_receiver/producer"
	"github.com/swarajroy/toll_calculator/types"
)

type LogMiddleware struct {
	next producer.DataProducer
}

func NewLogMiddleware(next producer.DataProducer) producer.DataProducer {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) Publish(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuId": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing data to kafka")
	}(time.Now())
	return l.next.Publish(data)
}
