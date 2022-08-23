// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package eta

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"

	"golang.org/x/sync/errgroup"
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
func (s *BatchService) GetEstimatedJobs(ctx context.Context, latLon job.LatLon, jobs []*job.Job) ([]*job.EstimatedJob, error) {
	if s.nWorkers < 1 {
		return nil, fmt.Errorf("at least one worker required")
	}

	err := validateLatLon(jobs)

	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	chunk := make(chan []*job.Job)

	g.Go(func() error {
		defer close(chunk)

		chunks := toChunks(jobs, s.chunkSize)

		for _, ch := range chunks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case chunk <- ch:
			}
		}

		return nil
	})

	estimatedJobs := make(chan []*job.EstimatedJob)

	workers := int32(s.nWorkers)
	for i := 0; i < s.nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one closes
				if atomic.AddInt32(&workers, -1) == 0 {
					close(estimatedJobs)
				}
			}()

			for ch := range chunk {
				estimatedChunk, err := EstimateJobs(ctx, s.googleAPIKey, latLon, ch)
				if err != nil {
					return fmt.Errorf("EstimateJobs error: %w", err)
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case estimatedJobs <- estimatedChunk:
				}
			}
			return nil
		})
	}

	eJobs := make([]*job.EstimatedJob, 0)
	g.Go(func() error {
		for ej := range estimatedJobs {
			eJobs = append(eJobs, ej...)
		}
		return nil
	})

	return eJobs, g.Wait()
}

func toChunks(jobs []*job.Job, chunkSize int) [][]*job.Job {
	var r [][]*job.Job
	for i := 0; i < len(jobs); i += chunkSize {
		end := i + chunkSize
		if end > len(jobs) {
			end = len(jobs)
		}
		r = append(r, jobs[i:end])
	}

	return r
}
