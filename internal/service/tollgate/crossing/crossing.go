package crossing

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
)

const table = "tollgate_crossing"

type Crossing struct {
	Crossing tollgate.Crossing
}

func (g *Crossing) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &g)
}

type TollgateCrossing struct {
	ID         dom.SUID `db:"id"`
	TollgateID dom.SUID `db:"tollgate_id"`
	DriverID   dom.SUID `db:"driver_id"`
	Crossing   Crossing `db:"crossing"`
	CreatedAt  dao.Time `db:"created_at"`
}

func (t *TollgateCrossing) Insert(ctx dom.Context) error {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.CreatedAt = now
	exec := dao.Insert(table).
		Set("id", t.ID).
		Set("tollgate_id", t.TollgateID).
		Set("driver_id", t.DriverID).
		Set("crossing", t.Crossing).
		Set("created_at", t.CreatedAt)
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
		Select("driver_id").
		Select("crossing").
		Select("created_at")
	for _, w := range wheres {
		query.Where(w.Expr, w.Args...)
	}
	query.OrderBy(orderBy...).
		Limit(limit)
	ary := make([]*TollgateCrossing, 0, limit)
	err := query.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c TollgateCrossing
			if err := rows.Scan(&c.ID, &c.TollgateID, &c.DriverID, &c.Crossing, &c.CreatedAt); err != nil {
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
