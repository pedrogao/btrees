package common

import (
	"unsafe"
)

type GoSlice struct {
	Ptr unsafe.Pointer
	Len int
	Cap int
}

//go:noescape
//go:linkname MemMove runtime.memmove
//goland:noinspection GoUnusedParameter
func MemMove(to unsafe.Pointer, from unsafe.Pointer, n uintptr)

// InsertAt insert element into slice at index
func InsertAt[T any](slice []T, s int, v T) []T {
	slice = append(slice[:s+1], slice[s:]...)
	slice[s] = v
	return slice
}

// RemoveAt remove element from slice at index
func RemoveAt[T any](slice []T, s int) T {
	ret := slice[s]
	slice = append(slice[:s], slice[s+1:]...)
	return ret
}
