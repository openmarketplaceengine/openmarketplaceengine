// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

import "github.com/openmarketplaceengine/openmarketplaceengine/dao"

type CarType uint8

const (
	CarTypeSedan = iota
	CarTypeHatch
	CarTypeTwoDoor
	CarTypeCabrio
	CarTypeSUV
	CarTypeVan
	CarTypeBus
	CarTypeTruck
	CarTypeCargoVan
	CarTypeLimo
)

type CarClass uint8

const (
	CarClassStandard = iota
	CarClassEconomy
	CarClassCompact
	CarClassExecutive
	CarClassLuxury
)

const vehicleTable = "vehicle"

// Vehicle represents car properties.
type Vehicle struct {
	ID         SUID     `db:"id"`
	Driver     SUID     `db:"driver"` // current driver's id, nullable
	Model      string   `db:"model"`
	Make       string   `db:"make"`
	Color      string   `db:"color"`
	Plate      string   `db:"plate"` // license plate
	Class      CarClass `db:"class"`
	Type       CarType  `db:"type"`
	Year       uint16   `db:"year"`
	Capacity   uint8    `db:"capacity"`   // passenger capacity
	CargoVol   uint8    `db:"cargovol"`   // cargo volume
	WheelChair uint8    `db:"wheelchair"` // wheelchair accommodation (0: none, 1: one, 2: two, ...)
	ChildSeats uint8    `db:"childseats"` // number of child seats
	Comment    string   `db:"comment"`
}

// Persist saves Vehicle to the database.
func (v *Vehicle) Persist(ctx Context) error {
	return dao.ExecTX(ctx, v.Insert())
}

//-----------------------------------------------------------------------------

func (v *Vehicle) Insert() dao.Executable {
	sql := dao.Insert(vehicleTable)
	sql.Set("id", v.ID)
	sql.SetNonZero("driver", v.Driver) // can be null
	sql.Set("model", v.Model)
	sql.Set("make", v.Make)
	sql.Set("color", v.Color)
	sql.Set("plate", v.Plate)
	sql.Set("class", v.Class)
	sql.Set("type", v.Type)
	sql.Set("year", v.Year)
	sql.Set("capacity", v.Capacity)
	sql.SetNonZero("cargovol", v.CargoVol)
	sql.Set("wheelchair", v.WheelChair)
	sql.Set("childseats", v.ChildSeats)
	sql.SetNonZero("comment", v.Comment)
	return sql
}
