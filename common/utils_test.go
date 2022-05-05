package common

import (
	"testing"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestMemmove(t *testing.T) {
	assert := assert.New(t)

	src := []byte{1, 2, 3, 4, 5, 6}
	dest := make([]byte, 10, 10)

	spew.Dump(src)
	spew.Dump(dest)
	srcp := (*GoSlice)(unsafe.Pointer(&src))
	destp := (*GoSlice)(unsafe.Pointer(&dest))

	MemMove(destp.Ptr, srcp.Ptr, unsafe.Sizeof(byte(0))*6)

	spew.Dump(src)
	spew.Dump(dest)
	assert.Equal(dest[0], uint8(1))
	assert.Equal(dest[1], uint8(2))
	assert.Equal(dest[4], uint8(5))
}
