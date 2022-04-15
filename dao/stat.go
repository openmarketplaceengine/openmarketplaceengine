package dao

import "github.com/openmarketplaceengine/openmarketplaceengine/stat"

func init() {
	group := stat.Group("database", "PostgreSQL database stats")
	group.Add("pool", "Connection pool stats", poolStat)
}

func poolStat(_ Context) (interface{}, error) {
	k := stat.GetIntKeyVal()
	s := DB().Stats()
	k.Add("MaxOpenConn", int64(s.MaxOpenConnections))
	k.Add("CurrOpenConn", int64(s.OpenConnections))
	k.Add("BusyOpenConn", int64(s.InUse))
	k.Add("IdleOpenConn", int64(s.Idle))
	return k, nil
}
