package collector

import (
	"github.com/matiue/GAgent/utils"
)

func GetNetworkUsage(interfaceName string) float64 {
	return utils.ReadNetworkUsage(interfaceName)
}
