package letheql

import (
	"time"

	"github.com/kuoss/lethe/config"
)

var now = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func init() {
	config.LoadConfig()
	config.GetConfig().Set("retention.time", "3h")
	config.GetConfig().Set("retention.size", "10m")
	config.SetNow(now)
	config.SetLimit(1000)
	//config.SetLogRoot(filepath.Join("data", "log"))
}
