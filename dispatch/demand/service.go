package demand

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate/store"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/geohash"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/geoqueue"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/jobstore"
	"os"
)

var apiKey = os.Getenv("OME_GOOGLE_API_KEY")

type Service struct {
	geoQueue      *geoqueue.GeoQueue
	estimateStore *store.EstimateStore
	jobStore      *jobstore.JobStore
}

func NewService() *Service {
	geoQueue := geoqueue.NewGeoQueue(dao.Reds.StoreClient)
	estimateStore := store.NewEstimateStore(dao.Reds.StoreClient)
	jobStore := jobstore.NewJobStore(dao.Reds.StoreClient)
	return &Service{
		geoQueue:      geoQueue,
		estimateStore: estimateStore,
		jobStore:      jobStore,
	}
}

func (s *Service) GetEstimates(ctx context.Context, areaKey string, from geoqueue.LatLon, radiusMeters int) ([]*estimate.Estimate, error) {
	nearByMembers, err := s.geoQueue.PeekMany(ctx, areaKey, from, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("peek many error: %w", err)
	}

	geoHash := geohash.ToGeoHash(from.Lat, from.Lon, geohash.Precision100)
	all, err := s.estimateStore.GetAll(ctx, geoHash, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("get all error: %w", err)
	}

	//todo some way to detect no need to call google
	if len(all) == len(nearByMembers) {
		return all, nil
	}

	jobs, err := s.jobStore.GetByIds(ctx, areaKey, toIds(nearByMembers)...)
	if err != nil {
		return nil, fmt.Errorf("get estimates error: %w", err)
	}

	estimates, err := estimate.Estimates(ctx, apiKey, estimate.LatLon{
		Lat: from.Lat,
		Lon: from.Lon,
	}, transform(jobs))

	if err != nil {
		return nil, fmt.Errorf("get estimates error: %w", err)
	}

	return estimates, nil
}

func transform(jobs []*jobstore.Job) []*estimate.Request {
	res := make([]*estimate.Request, len(jobs))

	for i, job := range jobs {
		res[i] = &estimate.Request{
			ID: job.ID,
			PickUp: estimate.LatLon{
				Lat: job.PickUp.Lat,
				Lon: job.PickUp.Lon,
			},
			DropOff: estimate.LatLon{
				Lat: job.DropOff.Lat,
				Lon: job.DropOff.Lon,
			},
		}
	}
	return res
}

func toIds(members []*geoqueue.Member) []string {
	res := make([]string, len(members))
	for i, member := range members {
		res[i] = member.ID
	}
	return res
}

func (s *Service) AddJob(ctx context.Context, areaKey string, job *Job) error {
	j := &jobstore.Job{
		ID: job.ID,
		PickUp: jobstore.LatLon{
			Lat: job.PickUp.Lat,
			Lon: job.PickUp.Lon,
		},
		DropOff: jobstore.LatLon{
			Lat: job.DropOff.Lat,
			Lon: job.DropOff.Lon,
		},
	}
	err := s.jobStore.StoreOne(ctx, areaKey, j)
	if err != nil {
		return err
	}

	err = s.geoQueue.Enqueue(ctx, areaKey, geoqueue.Member{
		ID: j.ID,
		PickUp: geoqueue.LatLon{
			Lat: j.PickUp.Lat,
			Lon: j.PickUp.Lon,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteJobs(ctx context.Context, areaKey string, ids ...string) error {

	for _, id := range ids {

		err := s.jobStore.Remove(ctx, areaKey, id)
		if err != nil {
			return err
		}
		_, err = s.geoQueue.Dequeue(ctx, areaKey, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetJob(ctx context.Context, areaKey string, id string) (*Job, error) {
	jobs, err := s.jobStore.GetByIds(ctx, areaKey, id)
	if err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return nil, nil
	}

	m, err := s.geoQueue.PeekOne(ctx, areaKey, id)

	if err != nil {
		return nil, err
	}

	return &Job{
		ID: id,
		PickUp: LatLon{
			Lat: m.PickUp.Lat,
			Lon: m.PickUp.Lon,
		},
		DropOff: LatLon{
			Lat: jobs[0].DropOff.Lat,
			Lon: jobs[0].DropOff.Lon,
		},
	}, nil
}
