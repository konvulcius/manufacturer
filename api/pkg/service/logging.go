package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/manufacturer/api/pkg/v1"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

// PostDoLogicWithManufacturer ...
func (l *loggingMiddleware) PostDoLogicWithManufacturer(ctx context.Context, input []*v1.ManufacturerInput) (err error) {
	defer func(begin time.Time) {
		_ = l.wrap(err).Log("method", "PostDoLogicWithManufacturer", "err", err, "took", time.Since(begin))
	}(time.Now())

	return l.next.PostDoLogicWithManufacturer(ctx, input)
}

func (l *loggingMiddleware) wrap(err error) log.Logger {
	if err != nil {
		return level.Error(l.logger)
	}

	return log.NewNopLogger()
}

// NewLoggingMiddleware ...
func NewLoggingMiddleware(next Service, logger log.Logger) Service {
	return &loggingMiddleware{
		logger: logger,
		next:   next,
	}
}
