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
