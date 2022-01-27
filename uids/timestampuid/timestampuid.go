package timestampuid

import (
	"fmt"
	"time"
)

func NewTimestampUid() string {
	return fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
}
