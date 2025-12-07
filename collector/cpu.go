package collector

import (
	"github.com/matiue/GAgent/utils"
)

func GetCPUUsage() float64 {
	return utils.ReadCPUUsage()
}
