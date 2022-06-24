package b2

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_node_add(t *testing.T) {
	assert := assert.New(t)

	n := newNode(5)
	assert.Equal(n.max, 5)
	assert.Equal(n.count, 0)
	assert.Equal(len(n.items), 5)
	assert.Equal(n.isLeaf(), true)
	assert.Equal(n.full(), false)

	n.add(1, "1")
	n.add(2, "2")
	n.add(3, "3")
	assert.Equal(n.getSize(), 3)
	n.add(4, "4")
	n2 := n.add(5, "5")
	assert.Equal(n2.getSize(), 3)
	assert.Equal(n2.items[0].key, 3)
	assert.Equal(n.getSize(), 2)
	assert.Equal(n.items[0].key, 1)
}

func TestBTree_Insert(t *testing.T) {
	assert := assert.New(t)

	bTree := NewBTree(3)
	assert.Equal(bTree.min, 3)
	assert.Equal(bTree.max, 6)
	assert.Equal(bTree.n, 0)

	ok := bTree.Insert(1, "1")
	assert.True(ok)

	ok = bTree.Insert(2, "2")
	assert.True(ok)
	ok = bTree.Insert(3, "3")
	assert.True(ok)
	ok = bTree.Insert(4, "4")
	assert.True(ok)
	ok = bTree.Insert(5, "5")
	assert.True(ok)
	ok = bTree.Insert(6, "6")
	assert.True(ok)

	val := bTree.Search(4)
	assert.Equal(val, "4")

	ok = bTree.Insert(7, "7")
	assert.True(ok)
	ok = bTree.Insert(8, "8")
	assert.True(ok)
	ok = bTree.Insert(9, "9")
	assert.True(ok)
}

func TestBTree_Delete(t *testing.T) {
	assert := assert.New(t)

	bTree := NewBTree(3)
	assert.Equal(bTree.min, 3)
	assert.Equal(bTree.max, 6)
	assert.Equal(bTree.n, 0)

	ok := bTree.Insert(1, "1")
	assert.True(ok)

	ok = bTree.Insert(2, "2")
	assert.True(ok)
	ok = bTree.Insert(3, "3")
	assert.True(ok)
	ok = bTree.Insert(4, "4")
	assert.True(ok)
	ok = bTree.Insert(5, "5")
	assert.True(ok)
	ok = bTree.Insert(6, "6")
	assert.True(ok)

	ok = bTree.Delete(4)
	assert.True(ok)

	ok = bTree.Delete(5)
	assert.True(ok)
}

func TestBTree_Delete2(t *testing.T) {
	assert := assert.New(t)

	bTree := NewBTree(3)
	assert.Equal(bTree.min, 3)
	assert.Equal(bTree.max, 6)
	assert.Equal(bTree.n, 0)

	count := 7

	for i := 1; i <= count; i++ {
		ok := bTree.Insert(i, strconv.Itoa(i))
		assert.True(ok)
	}
	// left rotation
	ok := bTree.Delete(1)
	assert.True(ok)
	ok = bTree.Delete(4)
	assert.True(ok)
}

func TestBTree_Delete3(t *testing.T) {
	assert := assert.New(t)

	bTree := NewBTree(2)
	assert.Equal(bTree.min, 2)
	assert.Equal(bTree.max, 4)
	assert.Equal(bTree.n, 0)

	count := 10

	for i := 1; i <= count; i++ {
		ok := bTree.Insert(i, strconv.Itoa(i))
		assert.True(ok)
	}
	// right rotation
	ok := bTree.Delete(6)
	assert.True(ok)
	ok = bTree.Delete(4)
	assert.True(ok)
}
