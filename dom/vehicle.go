// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

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

// Vehicle represents car properties.
type Vehicle struct {
	ID         SUID
	Driver     SUID // current driver's id, nullable
	Model      string
	Make       string
	Color      string
	Plate      string // license plate
	Class      CarClass
	Type       CarType
	Year       uint16
	Capacity   uint8 // passenger capacity
	CargoVol   uint8 // cargo volume
	WheelChair uint8 // wheelchair accommodation (0: none, 1: one, 2: two, ...)
	ChildSeats uint8 // number of child seats
	Comment    string
}
