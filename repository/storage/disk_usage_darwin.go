package storage

import "syscall"

// DiskUsage returns disk usage stats for the given path on macOS.
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
