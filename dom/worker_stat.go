package dom

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
)

func init() {
	wg := stat.Group("workers", "Workers stats")
	wg.Add("status", "Workers status info", workerStatus)
}

func workerStatus(ctx Context) (interface{}, error) {
	k := stat.GetIntKeyVal()
	sql := dao.NewSQL(
		`select status, count(status)
			from worker
			group by status
			order by 1`)
	err := sql.QueryEach(ctx, func(rows *dao.Rows) {
		var status, count int64
		if err := rows.Scan(&status, &count); err != nil {
			return
		}
		k.Add(WorkerStatus(status).String(), count)
	})
	if k.Len() > 0 {
		k.Add("total", k.Total())
	}
	return k, err
}
