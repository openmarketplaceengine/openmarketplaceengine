package dao

import (
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
)

func init() {
	group := stat.Group("database", "PostgreSQL database stats")
	group.Add("pool", "Connection pool stats", poolStat)
}

func poolStat(_ Context) (interface{}, error) {
	k := stat.GetIntKeyVal()
	s := DB().Stats()
	k.Add("open_conn_count", int64(s.OpenConnections))
	k.Add("open_conn_busy", int64(s.InUse))
	k.Add("open_conn_idle", int64(s.Idle))
	k.Add("conn_wait_count", s.WaitCount)
	k.Add("conn_wait_msecs", int64(s.WaitDuration/time.Millisecond))
	return k, nil
}
