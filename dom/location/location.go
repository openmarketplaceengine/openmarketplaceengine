package location

import (
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"gopkg.in/yaml.v2"
)

const table = "location"

type Location struct {
	Recnum    int64    `db:"recnum"` // auto-increasing record number
	Worker    dom.SUID `db:"worker"`
	Stamp     dom.Time `db:"stamp"`
	Longitude float64  `db:"longitude"`
	Latitude  float64  `db:"latitude"`
	Speed     int      `db:"speed"`
}

func QueryLast(ctx dom.Context, workerID dom.SUID) (loc *Location, has bool, err error) {
	loc = new(Location)
	sql := dao.From(table).
		Bind(loc).
		Where("worker = ?", workerID).
		OrderBy("stamp desc").
		Limit(1)
	has, err = sql.QueryOne(ctx)
	return
}

func QueryAll(ctx dom.Context, workerID dom.SUID, limit int) ([]*Location, error) {
	sql := dao.From(table).
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

// Insert inserts Location to the database.
func (l *Location) Insert(ctx dao.Context) error {
	executable := l.insert()
	return dao.ExecTX(ctx, executable)
}

func (l *Location) insert() dao.Executable {
	now := dao.Time{}
	now.Now()
	now.UTC()
	l.Stamp = now
	return dao.Insert(table).
		Set("worker", l.Worker).
		Set("stamp", l.Stamp).
		Set("longitude", l.Longitude).
		Set("latitude", l.Latitude).
		Set("speed", l.Speed)
}

func (l *Location) DumpYAML() {
	buf, err := yaml.Marshal(l)
	if err != nil {
		fmt.Printf("yaml.Marshal failed: %s\n", err)
		return
	}
	fmt.Printf("%s", buf)
}
