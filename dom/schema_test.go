package dom

import (
	"testing"

	"github.com/driverscooperative/geosrv/dao"
)

func TestDropAll(t *testing.T) {
	dao.AutoDrop(dropAll)
	dao.WillTest(t, "test")
}

func TestSchema(t *testing.T) {
	WillTest(t, "test", false)
}
