package dom

import (
	"fmt"

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

// Persist saves WorkerLocation to the database.
func (w *WorkerLocation) Persist(ctx dao.Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

//-----------------------------------------------------------------------------

func (w *WorkerLocation) Insert() dao.Executable {
	sql := dao.Insert(workerLocationTable)
	sql.Set("worker", w.Worker)
	sql.Set("stamp", w.Stamp)
	sql.Set("longitude", w.Longitude)
	sql.Set("latitude", w.Latitude)
	sql.Set("speed", w.Speed)
	sql.Returning("recnum").To(&w.Recnum)
	return sql
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
