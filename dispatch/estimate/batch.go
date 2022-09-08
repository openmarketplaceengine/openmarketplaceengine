// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package estimate

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
)

type BatchService struct {
	nWorkers     int
	chunkSize    int
	googleAPIKey string
}

func NewBatchService(googleAPIKey string, nWorkers int, chunkSize int) *BatchService {
	return &BatchService{
		nWorkers:     nWorkers,
		chunkSize:    chunkSize,
		googleAPIKey: googleAPIKey,
	}
}
func (s *BatchService) GetEstimates(ctx context.Context, latLon LatLon, requests []*Request) ([]*Estimate, error) {
	if s.nWorkers < 1 {
		return nil, fmt.Errorf("at least one worker required")
	}

	err := validateLatLon(requests)

	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	chunk := make(chan []*Request)

	g.Go(func() error {
		defer close(chunk)

		chunks := toChunks(requests, s.chunkSize)

		for _, ch := range chunks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case chunk <- ch:
			}
		}

		return nil
	})

	estimated := make(chan []*Estimate)

	workers := int32(s.nWorkers)
	for i := 0; i < s.nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one closes
				if atomic.AddInt32(&workers, -1) == 0 {
					close(estimated)
				}
			}()

			for ch := range chunk {
				estimatedChunk, err := Estimates(ctx, s.googleAPIKey, latLon, ch)
				if err != nil {
					return fmt.Errorf("Estimates error: %w", err)
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case estimated <- estimatedChunk:
				}
			}
			return nil
		})
	}

	estimates := make([]*Estimate, 0)
	g.Go(func() error {
		for e := range estimated {
			estimates = append(estimates, e...)
		}
		return nil
	})

	return estimates, g.Wait()
}

func toChunks(requests []*Request, chunkSize int) [][]*Request {
	var r [][]*Request
	for i := 0; i < len(requests); i += chunkSize {
		end := i + chunkSize
		if end > len(requests) {
			end = len(requests)
		}
		r = append(r, requests[i:end])
	}

	return r
}
