// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/brokers"
)

var _ brokers.MessagePublisher = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     brokers.MessagePublisher
}

// MetricsMiddleware instruments adapter by tracking request count and latency.
func MetricsMiddleware(svc brokers.MessagePublisher, counter metrics.Counter, latency metrics.Histogram) brokers.MessagePublisher {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm *metricsMiddleware) Publish(ctx context.Context, token string, msg mainflux.Message) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "publish").Add(1)
		mm.latency.With("method", "publish").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Publish(ctx, token, msg)
}
