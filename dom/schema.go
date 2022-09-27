package dom

import (
	"embed"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/dao"
)

const sqlPath = "schema/v1/"

//go:embed schema/v1/index.yaml
//go:embed schema/v1/*.sql
var sqlFsys embed.FS

//go:embed schema/v1/drop_all.sql
var dropAll dao.SQLExec

func Boot() error {
	if cfg.Pgdb().ReadyForStart() {
		dao.AutoExec(dao.NewFsysExec(sqlFsys, sqlPath, "index.yaml"))
	}
	return nil
}

func WillTest(t dao.Tester, schema string, drop bool) {
	if drop {
		dao.AutoDrop(dropAll)
	}
	_ = Boot()
	dao.WillTest(t, schema)
}
