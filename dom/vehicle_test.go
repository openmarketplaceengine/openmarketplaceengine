package dom

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestVehicle_Persist(t *testing.T) {
	WillTest(t, "test", true)
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		car := genVehicle()
		require.NoError(t, car.Persist(ctx))
	}
}

//-----------------------------------------------------------------------------

func genVehicle() *Vehicle {
	return &Vehicle{
		ID:         mockUUID("car"),
		Driver:     mockUUID("drv"),
		Model:      mockString("M1", "M2", "M3", "M4", "M5", "M6"),
		Make:       mockString("Toyota", "BMW", "Honda", "Ford", "Tesla", "GMC"),
		Color:      mockString("Black", "White", "Red", "Blue", "Navy"),
		Plate:      mockString("NYC-001", "NYC-002", "NYC-003", "NYC-004", "NYC-005"),
		Class:      CarClass(mockEnum(CarClassLuxury)),
		Type:       CarType(mockEnum(CarTypeLimo)),
		Year:       uint16(mockRange(1990, 2022)),
		Capacity:   uint8(mockRange(1, 4)),
		CargoVol:   uint8(mockRange(0, 180)),
		WheelChair: uint8(mockRange(0, 1)),
		ChildSeats: uint8(mockRange(0, 1)),
	}
}
