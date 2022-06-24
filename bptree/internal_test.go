package bptree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_internalNode_insert(t *testing.T) {
	assert := assert.New(t)
	// 2 个 key，3 个 pointer，以 key 作为 full，split 的标准
	n := newInternalNode(3)
	assert.Equal(n.kcs.Len(), 3)
	p1 := newLeafNode(2)
	p1.insert(1, "c")
	p1.insert(5, "a")
	n.insert(0, p1) // 0表示无

	p2 := newLeafNode(2)
	p2.insert(12, "b")
	n.insert(5, p2)

	assert.Equal(n.count, 2)
	assert.False(n.full())
}

func Test_internalNode_split(t *testing.T) {
	assert := assert.New(t)
	// 2 个 key，3 个 pointer，以 key 作为 full，split 的标准
	n := newInternalNode(3)
	assert.Equal(n.kcs.Len(), 3)

	p1 := newLeafNode(3)
	p1.insert(1, "1")
	n.insert(0, p1) // 0表示无

	p2 := newLeafNode(3)
	p2.insert(5, "5")
	n.insert(5, p2)

	assert.Equal(n.count, 2)
	assert.False(n.full())

	p3 := newLeafNode(3)
	p3.insert(12, "12")
	n.insert(12, p3)

	p4 := newLeafNode(3)
	p4.insert(18, "18")
	p4.insert(21, "21")
	n.insert(18, p4)
}

func Test_internalNode_full(t *testing.T) {
	assert := assert.New(t)

	n := newInternalNode(3)
	assert.Equal(n.kcs.Len(), 3)
	n.count = 1
	assert.Equal(n.halfFull(), true)
	n.count = 2
	assert.Equal(n.halfFull(), true)
	assert.Equal(n.full(), false)
	n.count = 3
	assert.Equal(n.full(), true)
	assert.Equal(n.getMaxSize(), 3)
	assert.Equal(n.getSize(), 3)
	assert.Equal(n.isLeaf(), false)
	n.resize(-1)
	assert.Equal(n.getSize(), 2)
	n.resize(2)
	assert.Equal(n.getSize(), 4)
}

func Test_internalNode_remove(t *testing.T) {
	assert := assert.New(t)

	n := newInternalNode(3)
	assert.Equal(n.kcs.Len(), 3)

	child := newLeafNode(3)
	child.insert(5, "a")
	child.insert(12, "b")
	child.insert(1, "c")
	child2 := newLeafNode(3)
	child3 := newLeafNode(3)

	n.insert(5, child3)
	n.insert(12, child2)
	n.insert(1, child)
	found, ok := n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, true)
	got := n.lookup(1)
	assert.Equal(got, child)
	n.remove(child)
	found, ok = n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, false)

	// n 5,12
	other := newInternalNode(3)
	assert.Equal(other.kcs.Len(), 3)
	n.moveLastToFrontOf(other)
	// n 5; other 12
	assert.Equal(n.getSize(), 1)
	assert.Equal(other.getSize(), 1)
	assert.Equal(n.kcs[0].key, 5)
	assert.Equal(other.kcs[0].key, 12)

	other.kcs[0].key = 13
	other.moveAllTo(n)
	// n 5, 13
	assert.Equal(n.getSize(), 2)
	assert.Equal(other.getSize(), 0)
	assert.Equal(n.kcs[0].key, 5)
	assert.Equal(n.kcs[1].key, 13)

	n.moveFirstToEndOf(other)
	// n 13; other 5
	assert.Equal(n.getSize(), 1)
	assert.Equal(other.getSize(), 1)
	assert.Equal(n.kcs[0].key, 13)
	assert.Equal(other.kcs[0].key, 5)
}
