package worker

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"gopkg.in/yaml.v2"
)

const workerLocationTable = "worker_location"

type Location struct {
	Recnum    int64    `db:"recnum"` // auto-increasing record number
	Worker    dom.SUID `db:"worker"`
	Stamp     dom.Time `db:"stamp"`
	Longitude float64  `db:"longitude"`
	Latitude  float64  `db:"latitude"`
	Speed     int      `db:"speed"`
}

//-----------------------------------------------------------------------------
// Getters
//-----------------------------------------------------------------------------

func LastLocation(ctx dom.Context, workerID dom.SUID) (loc *Location, has bool, err error) {
	loc = new(Location)
	sql := dao.From(workerLocationTable).
		Bind(loc).
		Where("worker = ?", workerID).
		OrderBy("stamp desc").
		Limit(1)
	has, err = sql.QueryOne(ctx)
	return
}

//-----------------------------------------------------------------------------

func ListLocations(ctx dom.Context, workerID dom.SUID, limit int) ([]*Location, error) {
	sql := dao.From(workerLocationTable).
		Select("recnum").
		Select("worker").
		Select("stamp").
		Select("longitude").
		Select("latitude").
		Select("speed").
		Where("worker = ?", workerID).
		OrderBy("stamp desc").
		Limit(limit)
	ary := make([]*Location, 0, limit)
	err := sql.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c Location
			if err := rows.Scan(
				&c.Recnum,
				&c.Worker,
				&c.Stamp,
				&c.Longitude,
				&c.Latitude,
				&c.Speed,
			); err != nil {
				return err
			}
			ary = append(ary, &c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ary, nil
}

//-----------------------------------------------------------------------------
// Setters
//-----------------------------------------------------------------------------

func AddLocation(ctx dom.Context, workerID dom.SUID, longitude float64, latitude float64, stamp time.Time, speed int) error {
	sql := dao.Insert(workerLocationTable).
		Set("worker", workerID).
		Set("stamp", stamp).
		Set("longitude", longitude).
		Set("latitude", latitude).
		Set("speed", speed)
	return dao.ExecTX(ctx, sql)
}

//-----------------------------------------------------------------------------

// Persist saves Location to the database.
func (l *Location) Persist(ctx dao.Context) error {
	return dao.ExecTX(ctx, l.Insert())
}

//-----------------------------------------------------------------------------

func (l *Location) Insert() dao.Executable {
	return dao.Insert(workerLocationTable).
		Set("worker", l.Worker).
		Set("stamp", l.Stamp).
		Set("longitude", l.Longitude).
		Set("latitude", l.Latitude).
		Set("speed", l.Speed)
}

//-----------------------------------------------------------------------------

func (l *Location) DumpYAML() {
	buf, err := yaml.Marshal(l)
	if err != nil {
		fmt.Printf("yaml.Marshal failed: %s\n", err)
		return
	}
	fmt.Printf("%s", buf)
}
