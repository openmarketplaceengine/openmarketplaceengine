package dispatch

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

var geoQueue = geoqueue.NewGeoQueue(dao.Reds.StoreClient)
var estimateStore = store.NewEstimateStore(dao.Reds.StoreClient)
var jobStore = jobstore.NewJobStore(dao.Reds.StoreClient)

var apiKey = os.Getenv("OME_GOOGLE_API_KEY")

func GetEstimates(ctx context.Context, areaKey string, from geoqueue.LatLon, radiusMeters int) ([]*estimate.Estimate, error) {
	nearByMembers, err := geoQueue.PeekMany(ctx, areaKey, from, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("peek many error: %w", err)
	}

	geoHash := geohash.ToGeoHash(from.Lat, from.Lon, geohash.Precision100)
	all, err := estimateStore.GetAll(ctx, geoHash, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("get all error: %w", err)
	}

	//todo some way to detect no need to call google
	if len(all) == len(nearByMembers) {
		return all, nil
	}

	jobs, err := jobStore.GetAll(ctx, areaKey, toIds(nearByMembers)...)
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

func AddJob(ctx context.Context, areaKey string, job *jobstore.Job) error {
	err := jobStore.StoreOne(ctx, areaKey, job)
	if err != nil {
		return err
	}

	err = geoQueue.Enqueue(ctx, areaKey, geoqueue.Member{
		ID: job.ID,
		PickUp: geoqueue.LatLon{
			Lat: job.PickUp.Lat,
			Lon: job.PickUp.Lon,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
