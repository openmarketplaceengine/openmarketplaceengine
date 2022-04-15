package stat

import "time"

var start = time.Now()

func init() {
	AddStat("_stamp", "Stat generation timestamp", stamp)
	AddStat("uptime", "Server uptime in seconds", uptime)
}

func stamp(_ Context) (interface{}, error) {
	return time.Now(), nil
}

func uptime(_ Context) (interface{}, error) {
	return time.Since(start), nil
}
