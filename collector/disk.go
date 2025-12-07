package collector

import (
	"github.com/matiue/GAgent/utils"
)

func GetDiskUsage(path string) float64 {
	return utils.ReadDiskUsage(path)
}
