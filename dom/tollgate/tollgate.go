package tollgate

import (
	"encoding/json"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
)

const table = "tollgate"

type BBoxes struct {
	BBoxes   []*detector.BBox
	Required int32
}

func (b *BBoxes) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &b)
}

type GateLine struct {
	Line *detector.Line
}

func (g *GateLine) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &g)
}

type Tollgate struct {
	ID       dom.SUID  `db:"id"`
	Name     string    `db:"name"`
	BBoxes   *BBoxes   `db:"b_boxes"`
	GateLine *GateLine `db:"gate_line"`
	Created  dao.Time  `db:"created"`
	Updated  dao.Time  `db:"updated"`
}

func (t *Tollgate) insert() dao.Executable {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.Created = now
	sql := dao.Insert(table).
		Set("id", t.ID).
		Set("name", t.Name).
		Set("b_boxes", t.BBoxes).
		Set("gate_line", t.GateLine).
		Set("created", t.Created)
	sql.IgnoreConflict()
	return sql
}

func (t *Tollgate) update() dao.Executable {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.Updated = now
	sql := dao.Update(table).
		Set("name", t.Name).
		Set("b_boxes", t.BBoxes).
		Set("gate_line", t.GateLine).
		Set("updated", t.Updated).
		Where("id = ?", t.ID)
	return sql
}

func (t *Tollgate) Insert(ctx dom.Context) error {
	executable := t.insert()
	return dao.ExecTX(ctx, executable)
}

func (t *Tollgate) Update(ctx dom.Context) error {
	executable := t.update()
	return dao.ExecTX(ctx, executable)
}

func (t *Tollgate) Upsert(ctx dom.Context) (dao.Result, dao.UpsertStatus, error) {
	return dao.Upsert(ctx, t.insert, t.update)
}

func QueryOne(ctx dom.Context, id dom.SUID) (*Tollgate, error) {
	var t Tollgate
	has, err := dao.From(table).
		Bind(&t).
		Where("id = ?", id).
		QueryOne(ctx)
	if has {
		return &t, nil
	}
	return nil, dao.WrapNoRows(err)
}

func QueryAll(ctx dom.Context, limit int) ([]*Tollgate, error) {
	query := dao.From(table).
		Select("id").
		Select("name").
		Select("b_boxes").
		Select("gate_line").
		Select("created").
		Select("updated").
		OrderBy("created desc").
		Limit(limit)
	ary := make([]*Tollgate, 0, limit)
	err := query.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c Tollgate
			if err := rows.Scan(&c.ID, &c.Name, &c.BBoxes, &c.GateLine, &c.Created, &c.Updated); err != nil {
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
