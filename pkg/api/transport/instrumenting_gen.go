//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	"context"
	"time"

	"github.com/fesyunoff/availability/pkg/api/service"
	"github.com/fesyunoff/availability/pkg/types"
	"github.com/go-kit/kit/metrics"
	prometheus2 "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type instrumentingOpts struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	namespace      string
	subsystem      string
}

type InstrumentingOption func(*instrumentingOpts)

func Namespace(v string) InstrumentingOption {
	return func(o *instrumentingOpts) {
		o.namespace = v
	}
}

func Subsystem(v string) InstrumentingOption {
	return func(o *instrumentingOpts) {
		o.subsystem = v
	}
}

func RequestLatency(requestLatency metrics.Histogram) InstrumentingOption {
	return func(o *instrumentingOpts) {
		o.requestLatency = requestLatency
	}
}

func RequestCount(requestCount metrics.Counter) InstrumentingOption {
	return func(o *instrumentingOpts) {
		o.requestCount = requestCount
	}
}

type ScraperRequestInstrumentingMiddleware struct {
	next service.ScraperRequest
	opts *instrumentingOpts
}

func (s *ScraperRequestInstrumentingMiddleware) GetAvailability(ctx context.Context, site string, id string) (string, error) {
	defer func(begin time.Time) {
		s.opts.requestCount.With("method", "GetAvailability").Add(1)
		s.opts.requestLatency.With("method", "GetAvailability").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.GetAvailability(ctx, site, id)
}

func (s *ScraperRequestInstrumentingMiddleware) GetResponceTime(ctx context.Context, limit string, id string) (string, error) {
	defer func(begin time.Time) {
		s.opts.requestCount.With("method", "GetResponceTime").Add(1)
		s.opts.requestLatency.With("method", "GetResponceTime").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.GetResponceTime(ctx, limit, id)
}

func (s *ScraperRequestInstrumentingMiddleware) GetStatistics(ctx context.Context, hours string, limit string, id string) ([]types.Stat, error) {
	defer func(begin time.Time) {
		s.opts.requestCount.With("method", "GetStatistics").Add(1)
		s.opts.requestLatency.With("method", "GetStatistics").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.next.GetStatistics(ctx, hours, limit, id)
}

func NewInstrumentingScraperRequestMiddleware(s service.ScraperRequest, opts ...InstrumentingOption) service.ScraperRequest {
	i := &ScraperRequestInstrumentingMiddleware{next: s, opts: &instrumentingOpts{}}
	for _, o := range opts {
		o(i.opts)
	}
	if i.opts.requestCount == nil {
		i.opts.requestCount = prometheus2.NewCounterFrom(prometheus.CounterOpts{
			Namespace: i.opts.namespace,
			Subsystem: i.opts.subsystem,
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"})

	}
	if i.opts.requestLatency == nil {
		i.opts.requestLatency = prometheus2.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: i.opts.namespace,
			Subsystem: i.opts.subsystem,
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"})

	}
	return i
}
