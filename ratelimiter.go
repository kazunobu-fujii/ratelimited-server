package rls

import (
	"context"
	"sort"

	"golang.org/x/time/rate"
)

// RateLimiter ...
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

// NewMultiLimiter ...
func NewMultiLimiter(limiters ...RateLimiter) *MultiLimiter {
	sort.Slice(limiters, func(lhs, rhs int) bool {
		return limiters[lhs].Limit() < limiters[rhs].Limit()
	})
	return &MultiLimiter{limiters: limiters}
}

// MultiLimiter ...
type MultiLimiter struct {
	limiters []RateLimiter
}

// Wait ...
func (l *MultiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Limit ...
func (l *MultiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}
