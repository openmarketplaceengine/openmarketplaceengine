package dom

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
)

var _rnd = rand.New(rand.NewSource(0))

func mockUUID(prefix ...string) string {
	if len(prefix) > 0 {
		prefix = append(prefix, dao.MockUUID())
		return strings.Join(prefix, "-")
	}
	return dao.MockUUID()
}

func realUUID(prefix ...string) string { //nolint
	if len(prefix) > 0 {
		prefix = append(prefix, dao.NewXid())
		return strings.Join(prefix, "-")
	}
	return dao.NewXid()
}

func numUUID(n int, prefix ...string) string { //nolint
	uuid := fmt.Sprintf("%08x", n)
	if len(prefix) > 0 {
		prefix = append(prefix, uuid)
		return strings.Join(prefix, "-")
	}
	return uuid
}

func mockStamp() (time Time) {
	time.Now()
	return
}

func mockSpeed() int {
	return _rnd.Intn(111)
}

func mockCoord() float64 {
	return _rnd.Float64() * 10_000
}

func mockEnum(last int) int {
	return _rnd.Intn(last + 1)
}

func mockString(list ...string) string {
	return list[_rnd.Intn(len(list))]
}

func mockRange(min, max int) int {
	return _rnd.Intn(max-min+1) + min
}

//-----------------------------------------------------------------------------

func dumpJSON(v interface{}) { //nolint:deadcode
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		println("dumpJSON failed:", err.Error())
		return
	}
	fmt.Printf("%s\n", buf)
}
