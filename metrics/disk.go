package metrics

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

type DiskUsage struct {
	MountPoint  string  `json:"mount_point"`
	Filesystem  string  `json:"filesystem"`
	TotalBytes  uint64  `json:"total_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

func GetDiskUsage(mount string) (DiskUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mount, &stat); err != nil {
		return DiskUsage{}, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	var usedPct float64
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100
	}
	return DiskUsage{
		MountPoint:  mount,
		TotalBytes:  total,
		FreeBytes:   free,
		UsedBytes:   used,
		UsedPercent: usedPct,
	}, nil
}

func GetAllDisks() ([]DiskUsage, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	disks := []DiskUsage{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		fsType := fields[2]
		mount := fields[1]

		if strings.HasPrefix(fsType, "tmpfs") ||
			strings.HasPrefix(fsType, "proc") ||
			strings.HasPrefix(fsType, "sysfs") ||
			strings.HasPrefix(fsType, "cgroup") ||
			strings.HasPrefix(fsType, "devpts") {
			continue
		}

		du, err := GetDiskUsage(mount)
		if err != nil {
			continue
		}
		du.Filesystem = fsType
		disks = append(disks, du)
	}
	return disks, scanner.Err()
}
