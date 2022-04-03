package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

const table = "tollgate"

type BBoxes struct {
	BBoxes   []*tollgate.BBox
	Required int
}

func (b *BBoxes) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &b)
}

type GateLine struct {
	Line tollgate.Line
}

func (g *GateLine) Scan(data interface{}) (err error) {
	return json.Unmarshal(data.([]byte), &g)
}

type Tollgate struct {
	ID        dom.SUID  `db:"id"`
	Name      string    `db:"name"`
	BBoxes    *BBoxes   `db:"b_boxes"`
	GateLine  *GateLine `db:"gate_line"`
	CreatedAt dao.Time  `db:"created_at"`
	UpdatedAt dao.Time  `db:"updated_at"`
}

func (t *Tollgate) Insert(ctx dom.Context) error {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.CreatedAt = now
	exec := dao.Insert(table).
		Set("id", t.ID).
		Set("name", t.Name).
		Set("b_boxes", t.BBoxes).
		Set("gate_line", t.GateLine).
		Set("created_at", t.CreatedAt)
	return dao.ExecTX(ctx, exec)
}

func (t *Tollgate) Update(ctx dom.Context) error {
	now := dao.Time{}
	now.Now()
	now.UTC()
	t.UpdatedAt = now
	exec := dao.Update(table).
		Set("name", t.Name).
		Set("b_boxes", t.BBoxes).
		Set("gate_line", t.GateLine).
		Set("updated_at", t.UpdatedAt).
		Where("id = ?", t.ID)
	return dao.ExecTX(ctx, exec)
}

func CreateIfNotExists(ctx dom.Context, tollgate *Tollgate) error {
	_, err := QueryOne(ctx, tollgate.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tollgate.Insert(ctx)
		}
		return fmt.Errorf("CreateIfNotExists, query tollgate error: %w", err)
	}
	return nil
}

func QueryOne(ctx dom.Context, id dom.SUID) (*Tollgate, error) {
	var t Tollgate
	err := dao.From(table).
		Bind(&t).
		Where("id = ?", id).
		QueryOne(ctx)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func QueryAll(ctx dom.Context, limit int) ([]*Tollgate, error) {
	query := dao.From(table).
		Select("id").
		Select("name").
		Select("b_boxes").
		Select("gate_line").
		Select("created_at").
		Select("updated_at").
		OrderBy("created_at desc").
		Limit(limit)
	ary := make([]*Tollgate, 0, limit)
	err := query.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var c Tollgate
			if err := rows.Scan(&c.ID, &c.Name, &c.BBoxes, &c.GateLine, &c.CreatedAt, &c.UpdatedAt); err != nil {
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
