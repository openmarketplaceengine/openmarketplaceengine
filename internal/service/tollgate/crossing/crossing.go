package crossing

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
)

const table = "tollgate_crossing"

type Crossing struct {
	Crossing detector.Crossing
}

func (g *Crossing) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &g)
}

type TollgateCrossing struct {
	ID         dom.SUID `db:"id"`
	TollgateID dom.SUID `db:"tollgate_id"`
	WorkerID   dom.SUID `db:"worker_id"`
	Crossing   Crossing `db:"crossing"`
	Created    dao.Time `db:"created"`
}

func (t *TollgateCrossing) Insert(ctx dom.Context) error {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.Created = now
	exec := dao.Insert(table).
		Set("id", t.ID).
		Set("tollgate_id", t.TollgateID).
		Set("worker_id", t.WorkerID).
		Set("crossing", t.Crossing).
		Set("created", t.Created)
	return dao.ExecTX(ctx, exec)
}

type Where struct {
	Expr string
	Args []interface{}
}

func QueryBy(ctx dom.Context, wheres []Where, orderBy []string, limit int) ([]*TollgateCrossing, error) {
	query := dao.From(table).
		Select("id").
		Select("tollgate_id").
		Select("worker_id").
		Select("crossing").
		Select("created")
	for _, w := range wheres {
		query.Where(w.Expr, w.Args...)
	}
	query.OrderBy(orderBy...).
		Limit(limit)
	ary := make([]*TollgateCrossing, 0, limit)
	err := query.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c TollgateCrossing
			if err := rows.Scan(&c.ID, &c.TollgateID, &c.WorkerID, &c.Crossing, &c.Created); err != nil {
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

func deleteAll(ctx dom.Context) error {
	del := dao.Delete(table)
	return dao.ExecTX(ctx, del)
}
