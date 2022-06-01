package detector

import (
	"context"
	"encoding/json"
	"fmt"
)

// Alg represents algorithm.
type Alg string

const (
	LineAlg   Alg = "line"
	BBoxAlg   Alg = "bbox"
	VectorAlg Alg = "vector"
)

type Handler func(ctx context.Context, c *Crossing) error

// Crossing represents detected fact of passing through the tollgate by WorkerID.
type Crossing struct {
	ID         string    `json:"id"`
	TollgateID string    `json:"tollgate_id"`
	WorkerID   string    `json:"worker_id"`
	Movement   *Movement `json:"movement"`
	Direction  Direction `json:"direction"`
	Alg        Alg       `json:"alg"`
}

func (c Crossing) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("json mardhar error: %s", err)
	}
	return string(bytes)
}

// Movement represents a moving from one Location to another.
type Movement struct {
	From *Location `json:"from"`
	To   *Location `json:"to"`
}

// Location is longitude, latitude corresponding to linear algebra X, Y axis.
type Location struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

//Direction to North, South, East or West in form of N, S, E, W, NE, NW, SE, SW.
type Direction string

// Direction of movement
// When moving to North - latitude increases until 90.
// When moving to South - latitude decreases until -90.
// When moving to East - longitude increases until 180.
// When moving to West - longitude decreases until -180.
// If 90/180 limit crossed it jumps to -90/-180 and vice versa.
// Movement represents a moving subject
// returns Direction in form of N, S, E, W, NE, NW, SE, SW.
func (m *Movement) Direction() Direction {
	fromX := m.From.Longitude
	fromY := m.From.Latitude
	toX := m.To.Longitude
	toY := m.To.Latitude

	var pole string
	var side string
	if fromY < toY {
		pole = "N"
	}
	if fromY > toY {
		pole = "S"
	}

	if fromX < toX {
		side = "E"
	}
	if fromX > toX {
		side = "W"
	}

	return Direction(fmt.Sprintf("%s%s", pole, side))
}

type Tollgate struct {
	ID             string
	Line           *Line
	BBoxes         []*BBox
	BBoxesRequired int32
}

type Detector struct {
	tollgates []*Tollgate
	storage   Storage
}

func NewDetector(tollgates []*Tollgate, storage Storage) *Detector {
	return &Detector{
		tollgates: tollgates,
		storage:   storage,
	}
}

func NewDetectorNoOp() *Detector {
	return &Detector{
		tollgates: []*Tollgate{},
		storage:   nil,
	}
}

// DetectCrossing detects if subject Movement has travelled through the tollgate.
func (d *Detector) DetectCrossing(ctx context.Context, workerID string, movement *Movement, handlers ...Handler) (*Crossing, error) {
	for _, t := range d.tollgates {
		if t.Line != nil {
			crossing := detectCrossingVector(t.ID, workerID, t.Line, movement)
			if crossing != nil {
				for _, handler := range handlers {
					innerErr := handler(ctx, crossing)
					if innerErr != nil {
						return nil, fmt.Errorf("handler error: %s", innerErr)
					}
				}
				return crossing, nil
			}
		}

		if len(t.BBoxes) > 0 {
			crossing, err := detectCrossingBBox(ctx, d.storage, t.ID, workerID, t.BBoxes, t.BBoxesRequired, movement)
			if err != nil {
				return nil, err
			}
			if crossing != nil {
				for _, handler := range handlers {
					innerErr := handler(ctx, crossing)
					if innerErr != nil {
						return nil, fmt.Errorf("handler error: %s", err)
					}
				}
				return crossing, nil
			}
		}
	}
	return nil, nil
}
