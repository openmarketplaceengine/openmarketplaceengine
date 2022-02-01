package timestampuid

import (
	"fmt"
	"time"
)

func NewTimestampUID() string {
	return fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
}
