package worker

import (
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

func randStatus() Status {
	return Status(_rnd.Intn(Disabled + 1))
}

func randInt(n int) int {
	return _rnd.Intn(n)
}

func randFirstName() string {
	n := _rnd.Intn(len(names))
	return names[n][0]
}

func randLastName() string {
	n := _rnd.Intn(len(names))
	return names[n][1]
}

var names = [][2]string{
	{"Jasmine", "Young"},
	{"Samantha", "Thomson"},
	{"Joseph", "Berry"},
	{"Leah", "Thomson"},
	{"Christopher", "Slater"},
	{"Carl", "Paterson"},
	{"Joan", "Skinner"},
	{"Caroline", "Parr"},
	{"Carl", "Peters"},
	{"Kimberly", "Watson"},
	{"Adam", "Bell"},
	{"Matt", "Clarkson"},
	{"Jan", "Jones"},
	{"Isaac", "Young"},
	{"Frank", "Ball"},
	{"Steven", "Baker"},
	{"Colin", "Churchill"},
	{"Jane", "James"},
	{"Warren", "Piper"},
	{"Oliver", "Avery"},
}
