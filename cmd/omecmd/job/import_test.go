// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

const (
	wkbSRID uint32 = 0x20000000 //nolint
)

/*
func TestGeoParse(t *testing.T) {
	geos := []string{
		"0101000020E6100000FE0C6FD6E07A52C02A00C633685C4440",
		"0101000020E6100000A1754309697B52C073E26190AA544440",
		// "0101000000000000000000F03F000000000000F03",
	}
	for i := range geos {
		val, err := hex.DecodeString(geos[i])
		require.NoError(t, err)
		ptyp := binary.LittleEndian.Uint32(val[1:])
		if ptyp&wkbSRID != 0 {
			ptyp = ptyp & ^wkbSRID
			t.Logf("%d: has SRID", i)
		}
		srid := binary.LittleEndian.Uint32(val[5:])
		pt1 := math.Float64frombits(binary.LittleEndian.Uint64(val[9:]))
		pt2 := math.Float64frombits(binary.LittleEndian.Uint64(val[17:]))
		t.Logf("%v\nlen: %d: [0]=%d, type: %d, srid: %d, pt1: %f, pt2: %f", val, len(val), val[0], ptyp, srid, pt1, pt2)
	}
}
*/
