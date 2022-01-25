package uids

import (
	"fmt"
	"time"
)

func GetTimestampUID() string {
	return fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
}
