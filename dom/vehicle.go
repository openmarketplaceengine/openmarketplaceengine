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
	return dao.Insert(vehicleTable).
		Set("id", v.ID).
		SetNonZero("driver", v.Driver).
		Set("model", v.Model).
		Set("make", v.Make).
		Set("color", v.Color).
		Set("plate", v.Plate).
		Set("class", v.Class).
		Set("type", v.Type).
		Set("year", v.Year).
		Set("capacity", v.Capacity).
		SetNonZero("cargovol", v.CargoVol).
		Set("wheelchair", v.WheelChair).
		Set("childseats", v.ChildSeats).
		SetNonZero("comment", v.Comment)
}
