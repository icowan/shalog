package post

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) Get(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "detail").Add(1)
		s.requestLatency.With("method", "detail").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Get(ctx, id)
}

func (s *instrumentingService) List(ctx context.Context, order, by, category string, pageSize, offset int) (rs []map[string]interface{}, count int64, other map[string]interface{}, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list").Add(1)
		s.requestLatency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.List(ctx, order, by, category, pageSize, offset)
}
