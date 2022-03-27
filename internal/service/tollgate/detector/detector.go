package detector

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/bbox"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/line"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/repository"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type Detector struct {
	tollgates []tollgate.Tollgate
}

func NewDetector() *Detector {
	return &Detector{
		tollgates: make([]tollgate.Tollgate, 0),
	}
}

func (d *Detector) AddTollgate(tollgate tollgate.Tollgate) {
	d.tollgates = append(d.tollgates, tollgate)
}

// LoadTollgates - loads tolls information from database.
// Consider encapsulating it in NewDetector.
func (d *Detector) LoadTollgates() error {
	tolls, err := repository.FindAll()
	if err != nil {
		return err
	}

	client := redisClient.NewStoreClient()

	storage := bbox.NewStorage(client)

	for _, t := range tolls.BboxTollgates {
		bboxT, innerErr := bbox.NewTollgate(t.ID, transform(t.Boxes), t.BoxesRequired, storage)
		if innerErr != nil {
			return nil
		}
		d.tollgates = append(d.tollgates, bboxT)
	}

	for _, t := range tolls.LineTollgates {
		lineT := line.NewTollgate(
			t.ID,
			&tollgate.LocationXY{
				LongitudeX: t.LonMin,
				LatitudeY:  t.LatMin,
			},
			&tollgate.LocationXY{
				LongitudeX: t.LonMax,
				LatitudeY:  t.LatMax,
			})
		d.tollgates = append(d.tollgates, lineT)
	}
	return nil
}

func transform(boxes []repository.Box) (bBoxes []*bbox.BBox) {
	for _, b := range boxes {
		bBoxes = append(bBoxes, &bbox.BBox{
			Left:   b.LonMin,
			Bottom: b.LatMin,
			Right:  b.LonMax,
			Top:    b.LatMax,
		})
	}
	return
}

func (d *Detector) DetectTollgateCrossing(ctx context.Context, movement *tollgate.Movement) *tollgate.Crossing {
	for _, t := range d.tollgates {
		crossing, err := t.DetectCrossing(ctx, movement)
		if err != nil {
			log.Errorf("detect tollgate crossing error: %q", err)
			continue
		}
		if crossing != nil {
			return crossing
		}
	}
	return nil
}
