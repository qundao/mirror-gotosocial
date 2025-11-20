package disk

import (
	"syscall"
	"time"
)

func modtime(sys *syscall.Stat_t) time.Time {
	return time.Unix(sys.Mtim.Unix())
}
