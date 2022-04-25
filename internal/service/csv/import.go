package csv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/crossing"
)

type Importer struct {
	tracker *location.Tracker
}

func NewImporter(tracker *location.Tracker) *Importer {
	return &Importer{
		tracker: tracker,
	}
}

func (i *Importer) Import(ctx context.Context, csvFile string) ([]*crossing.TollgateCrossing, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %w", csvFile, err)
	}
	s := bufio.NewScanner(file)

	var line int
	var crossings []*crossing.TollgateCrossing
	start := time.Now()
	for {
		scan := s.Scan()
		line++
		if !scan {
			err := s.Err()
			if err != nil {
				return nil, fmt.Errorf("scan error:%w", err)
			}
			fmt.Printf("Detected %v crossings in %s\n", len(crossings), time.Since(start))
			return crossings, nil
		}
		text := s.Text()
		isHeader := strings.Contains(text, "driver_id")
		if isHeader {
			continue
		}
		l, err := Parse(text)
		if err != nil {
			return nil, fmt.Errorf("line %v parse error:%w", line, err)
		}

		x, err := i.tracker.TrackLocation(ctx, l.DriverID, l.Lon, l.Lat)
		if err != nil {
			return nil, fmt.Errorf("track location error:%w", err)
		}
		if x != nil {
			crossings = append(crossings, x)
		}
	}
}
