package common

import (
	"fmt"
	"syscall"
)

func MmapAnonymous(size int) ([]byte, error) {
	// MAP_ANON：会忽略参数 fd，映射区不与任何文件关联，而且映射区域无法和其他进程共享
	prot, flags, fd, offset := syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE, -1, 0

	mm, err := syscall.Mmap(fd, int64(offset), size, prot, flags)
	if err != nil {
		return nil, fmt.Errorf("mmap err: %s", err)
	}

	return mm, nil
}

func Unmap(b []byte) error {
	return syscall.Munmap(b)
}
