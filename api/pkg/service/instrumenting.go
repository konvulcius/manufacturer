package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/manufacturer/api/pkg/v1"
)

type instrumentingMiddleware struct {
	next       Service
	reqMetrics metrics.Histogram
}

// PostDoLogicWithManufacturer ...
func (i *instrumentingMiddleware) PostDoLogicWithManufacturer(ctx context.Context, input []*v1.ManufacturerInput) (err error) {
	defer i.recordMetrics("PostDoLogicWithManufacturer", time.Now(), err)

	return i.next.PostDoLogicWithManufacturer(ctx, input)
}

func (i *instrumentingMiddleware) recordMetrics(method string, startTime time.Time, err error) {
	labels := []string{
		"method", method,
		"error", strconv.FormatBool(err != nil),
	}
	i.reqMetrics.With(labels...).Observe(time.Since(startTime).Seconds())
}

// NewInstrumentingMiddleware ...
func NewInstrumentingMiddleware(next Service, prefix string) Service {
	return &instrumentingMiddleware{
		next: next,
		reqMetrics: kitprometheus.NewHistogramFrom(
			prometheus.HistogramOpts{
				Name: prefix + "_requests",
				Help: "Requests Info",
				Buckets: []float64{
					0.001, 0.002, 0.003,
					0.004, 0.005, 0.010,
					0.015, 0.030, 0.050,
					0.100, 0.3, 0.5,
					1, 2, 5,
				},
			},
			[]string{"method", "error"},
		),
	}
}
