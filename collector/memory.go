package collector

import (
	"github.com/matiue/GAgent/utils"
)

func GetMemoryUsage() float64 {
	return utils.ReadMemoryUsage()
}
