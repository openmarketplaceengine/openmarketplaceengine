package dom

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"gopkg.in/yaml.v2"
)

const workerLocationTable = "worker_location"

type WorkerLocation struct {
	Recnum    int64   `db:"recnum"` // auto-increasing record number
	Worker    SUID    `db:"worker"`
	Stamp     Time    `db:"stamp"`
	Longitude float64 `db:"longitude"`
	Latitude  float64 `db:"latitude"`
	Speed     int     `db:"speed"`
}

//-----------------------------------------------------------------------------
// Getters
//-----------------------------------------------------------------------------

func LastWorkerLocation(ctx Context, workerID SUID) (loc Coord, err error) {
	sql := dao.From(workerLocationTable).
		Select("longitude").To(&loc.Longitude).
		Select("latitude").To(&loc.Latitude).
		Where("worker = ?", workerID).
		OrderBy("stamp desc").
		Limit(1)
	err = sql.QueryOne(ctx)
	return
}

//-----------------------------------------------------------------------------

func ListWorkerLocation(ctx Context, workerID SUID, limit int) ([]Coord, error) {
	sql := dao.From(workerLocationTable).
		Select("longitude").
		Select("latitude").
		Where("worker = ?", workerID).
		OrderBy("stamp desc").
		Limit(limit)
	ary := make([]Coord, 0, limit)
	err := sql.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c Coord
			if err := rows.Scan(&c.Longitude, &c.Latitude); err != nil {
				return err
			}
			ary = append(ary, c)
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

func AddWorkerLocation(ctx Context, workerID SUID, lon float64, lat float64, stamp time.Time, speed int) error {
	sql := dao.Insert(workerLocationTable).
		Set("worker", workerID).
		Set("stamp", stamp).
		Set("longitude", lon).
		Set("latitude", lat).
		Set("speed", speed)
	return dao.ExecTX(ctx, sql)
}

//-----------------------------------------------------------------------------

// Persist saves WorkerLocation to the database.
func (w *WorkerLocation) Persist(ctx dao.Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

//-----------------------------------------------------------------------------

func (w *WorkerLocation) Insert() dao.Executable {
	return dao.Insert(workerLocationTable).
		Set("worker", w.Worker).
		Set("stamp", w.Stamp).
		Set("longitude", w.Longitude).
		Set("latitude", w.Latitude).
		Set("speed", w.Speed)
}

//-----------------------------------------------------------------------------

func (w *WorkerLocation) DumpYAML() {
	buf, err := yaml.Marshal(w)
	if err != nil {
		fmt.Printf("yaml.Marshal failed: %s\n", err)
		return
	}
	fmt.Printf("%s", buf)
}
