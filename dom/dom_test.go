package dom

import (
	"fmt"
	"math/rand"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
)

var _rnd = rand.New(rand.NewSource(0))

func mockUUID() string {
	return dao.MockUUID()
}

func numUUID(n int) string { //nolint
	return fmt.Sprintf("%08x", n)
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
