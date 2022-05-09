package worker

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
)

func init() {
	wg := stat.App().Group("workers", "Workers stats")
	wg.Add("status", "Workers by status", workersByStatus)
}

func workersByStatus(ctx dom.Context) (interface{}, error) {
	//goland:noinspection GrazieInspection
	const query = `
		select a.name, count(worker.status)
		from
		    worker_status_enum a
		left join
		    worker on a.id = worker.status
		group by
		    a.id
		order by a.id
`
	k := stat.GetIntKeyVal()
	sql := dao.NewSQL(query)
	err := sql.QueryRows(ctx, func(rows *dao.Rows) error {
		return k.Scan(rows)
	})
	if k.Len() > 0 {
		k.Add("total", k.Total())
	}
	return k, err
}
