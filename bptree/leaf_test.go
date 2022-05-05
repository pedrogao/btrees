package bptree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_leafNode_find(t *testing.T) {
	assert := assert.New(t)

	n := newLeafNode(nil, 2)
	assert.Equal(n.kvs.Len(), 2)
	found, ok := n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, false)

	n.kvs[0] = kv{
		key:   5,
		value: "a",
	}
	n.kvs[1] = kv{
		key:   12,
		value: "b",
	}
	n.count = 2
	found, ok = n.find(5)
	assert.Equal(found, 0)
	assert.Equal(ok, true)
	found, ok = n.find(12)
	assert.Equal(found, 1)
	assert.Equal(ok, true)
}

func Test_leafNode_insert(t *testing.T) {
	assert := assert.New(t)

	n := newLeafNode(nil, 3)
	assert.Equal(n.kvs.Len(), 3)
	found, ok := n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, false)

	n.insert(5, "a")
	n.insert(12, "b")
	n.insert(1, "c")
	found, ok = n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, true)
}

func Test_leafNode_split(t *testing.T) {
	assert := assert.New(t)

	n := newLeafNode(nil, 4)
	assert.Equal(n.kvs.Len(), 4)
	n.kvs[0] = kv{
		key:   1,
		value: "c",
	}
	n.kvs[1] = kv{
		key:   5,
		value: "a",
	}
	n.kvs[2] = kv{
		key:   12,
		value: "b",
	}
	n.kvs[3] = kv{
		key:   20,
		value: "d",
	}
	n.count = 4
	next := n.split()
	assert.Equal(n.count, 2)
	assert.Equal(next.count, 2)
	assert.Equal(n.kvs[0].key, 1)
	assert.Equal(n.kvs[0].value, "c")
	assert.Equal(n.kvs[1].key, 5)
	assert.Equal(n.kvs[1].value, "a")
	assert.Equal(next.kvs[0].key, 12)
	assert.Equal(next.kvs[0].value, "b")
	assert.Equal(next.kvs[1].key, 20)
	assert.Equal(next.kvs[1].value, "d")
}

func Test_leafNode_full(t *testing.T) {
	assert := assert.New(t)

	n := newLeafNode(nil, 3)
	assert.Equal(n.kvs.Len(), 3)
	n.count = 1
	assert.Equal(n.halfFull(), true)
	n.count = 2
	assert.Equal(n.halfFull(), true)
	assert.Equal(n.full(), false)
	n.count = 3
	assert.Equal(n.full(), true)
	assert.Equal(n.getMax(), 3)
	assert.Equal(n.getMid(), 1)
	assert.Equal(n.getSize(), 3)
	assert.Equal(n.isLeaf(), true)
	n.resize(-1)
	assert.Equal(n.getSize(), 2)
	n.resize(2)
	assert.Equal(n.getSize(), 4)
}

func Test_leafNode_remove(t *testing.T) {
	assert := assert.New(t)

	n := newLeafNode(nil, 3)
	assert.Equal(n.kvs.Len(), 3)
	n.insert(5, "a")
	n.insert(12, "b")
	n.insert(1, "c")
	found, ok := n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, true)
	removed := n.remove(1)
	assert.Equal(removed, true)
	found, ok = n.find(1)
	assert.Equal(found, 0)
	assert.Equal(ok, false)

	// n 5,12
	other := newLeafNode(nil, 3)
	assert.Equal(other.kvs.Len(), 3)
	n.moveLastToFrontOf(other)
	// n 5; other 12
	assert.Equal(n.getSize(), 1)
	assert.Equal(other.getSize(), 1)
	assert.Equal(n.kvs[0].key, 5)
	assert.Equal(other.kvs[0].key, 12)

	other.kvs[0].key = 13
	other.moveAllTo(n)
	// n 5, 13
	assert.Equal(n.getSize(), 2)
	assert.Equal(other.getSize(), 0)
	assert.Equal(n.kvs[0].key, 5)
	assert.Equal(n.kvs[1].key, 13)

	n.moveFirstToEndOf(other)
	// n 13; other 5
	assert.Equal(n.getSize(), 1)
	assert.Equal(other.getSize(), 1)
	assert.Equal(n.kvs[0].key, 13)
	assert.Equal(other.kvs[0].key, 5)
}
