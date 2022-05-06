package bptree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func verifyTree(b *BPTree, count int, t *testing.T) {
	verifyRoot(b, t)

	for i := 0; i < b.root.(*internalNode).count; i++ {
		verifyNode(b.root.(*internalNode).kcs[i].child, b.root.(*internalNode), t)
	}

	leftMost := findLeftMost(b.root)

	verifyLeaf(leftMost, count, t)
}

// min child: 1
// max child: MaxKC
func verifyRoot(b *BPTree, t *testing.T) {
	if b.Empty() {
		t.Logf("empty tree")
		return
	}
	if b.root.parent() != nil {
		t.Errorf("root.parent: want = nil, got = %p", b.root.parent())
	}

	if b.root.getSize() < 1 {
		t.Errorf("root.min.child: want >=1, got = %d", b.root.(*internalNode).count)
	}

	if b.root.getSize() > b.root.getMaxSize() {
		t.Errorf("root.max.child: want <= %d, got = %d", b.root.(*internalNode).max, b.root.(*internalNode).count)
	}
}

func verifyNode(n node, parent *internalNode, t *testing.T) {
	switch nn := n.(type) {
	case *internalNode:
		if !nn.halfFull() {
			t.Errorf("internal.min.child: want >= %d, got = %d", nn.max/2, nn.count)
		}

		if nn.getSize() > nn.getMaxSize() {
			t.Errorf("internal.max.child: want <= %d, got = %d", nn.max, nn.count)
		}

		if nn.parent() != parent {
			t.Errorf("internal.parent: want = %p, got = %p", parent, nn.parent())
		}

		for i := 0; i < nn.count; i++ {
			if i > 0 {
				if nn.kcs[i].key < nn.kcs[i-1].key {
					t.Errorf("right = %d must bigger than left = %d", nn.kcs[i].key, nn.kcs[i-1].key)
				}
			}
			verifyNode(nn.kcs[i].child, nn, t)
		}

	case *leafNode:
		if nn.parent() != parent {
			t.Errorf("leaf.parent: want = %p, got = %p", parent, nn.parent())
		}

		if !nn.halfFull() {
			t.Errorf("leaf.min.child: want >= %d, got = %d", nn.max/2, nn.count)
		}

		if nn.getSize() > nn.getMaxSize() {
			t.Errorf("leaf.max.child: want <= %d, got = %d", nn.max, nn.count)
		}

		for i := 0; i < nn.count; i++ {
			if i > 0 {
				if nn.kvs[i].key < nn.kvs[i-1].key {
					t.Errorf("right = %d must bigger than left = %d", nn.kvs[i].key, nn.kvs[i-1].key)
				}
			}
		}
	}
}

func verifyLeaf(leftMost *leafNode, count int, t *testing.T) {
	curr := leftMost
	last := 0
	c := 0

	for curr != nil {
		for i := 0; i < curr.count; i++ {
			key := curr.kvs[i].key

			if key <= last {
				t.Errorf("leaf.sort.key: want > %d, got = %d", last, key)
			}
			last = key
			c++
		}
		curr = curr.next
	}

	if c != count {
		t.Errorf("leaf.count: want = %d, got = %d", count, c)
	}
}

func findLeftMost(n node) *leafNode {
	switch nn := n.(type) {
	case *internalNode:
		return findLeftMost(nn.kcs[0].child)
	case *leafNode:
		return nn
	default:
		panic("unknown node type")
	}
}

func TestBTree_Insert1(t1 *testing.T) {
	keys := []int{1, 5, 12, 18, 21, 22, 23}
	bt := NewBPTree(MaxInternal(3), MaxLeaf(3))

	for _, key := range keys {
		bt.Insert(key, fmt.Sprintf("%d", key))
		//bt.printTree()
	}
}

func TestBTree_Insert2(t1 *testing.T) {
	keys := []int{1, 5, 12, 18, 21, 22, 23}
	bt := NewBPTree(MaxInternal(6), MaxLeaf(6))

	for _, key := range keys {
		bt.Insert(key, fmt.Sprintf("%d", key))
		bt.printTree()
	}
}

func TestBTree_Search1(t1 *testing.T) {
	assert := assert.New(t1)
	bt := NewBPTree(MaxInternal(10), MaxLeaf(10))
	count := 1000

	for i := 1; i <= count; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	for i := 1; i <= count; i++ {
		got, ok := bt.Search(i)
		assert.True(ok, "expect=%d, but got=%s", i, got)
	}

	//graph := bt.Graph()
	//assert.NotEqual(graph, "")
	//err := ioutil.WriteFile("test.dot", []byte(graph), 0666)
	//assert.Nil(err)
}

func TestBTree_Search2(t1 *testing.T) {
	assert := assert.New(t1)
	bt := NewBPTree()
	count := 100000

	for i := 1; i <= count; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	for i := 1; i <= count; i++ {
		got, ok := bt.Search(i)
		assert.True(ok, "expect=%d, but got=%s", i, got)
	}

	verifyTree(bt, count, t1)
}

func TestBTree_Search3(t1 *testing.T) {
	assert := assert.New(t1)
	bt := NewBPTree()
	count := 100000

	for i := 1; i <= count; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	for i := 1; i <= count; i++ {
		got, ok := bt.Search(i)
		assert.True(ok, "expect=%d, but got=%s", i, got)
	}

	for i := 1; i <= count; i++ {
		bt.Delete(i)
	}

	assert.True(bt.Empty())

	for i := 1; i <= count; i++ {
		_, ok := bt.Search(i)
		assert.False(ok)
	}
}

func TestBTree_Delete2(t1 *testing.T) {
	//assert := assert.New(t1)
	bt := NewBPTree(MaxInternal(10), MaxLeaf(10))
	count := 100

	for i := 1; i <= count; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	for i := 1; i <= count; i++ {
		bt.Delete(i)
	}

	//graph := bt.Graph()
	//assert.NotEqual(graph, "")
	//err := ioutil.WriteFile("test.dot", []byte(graph), 0666)
	//assert.Nil(err)
}

func TestBTree_Delete1(t *testing.T) {
	keys := []int{1, 5, 12, 18, 21, 22, 23}
	bt := NewBPTree(MaxInternal(6), MaxLeaf(6))

	for _, key := range keys {
		bt.Insert(key, fmt.Sprintf("%d", key))
		bt.printTree()
	}

	t.Logf("delete begin")

	for _, key := range keys {
		bt.Delete(key)
		bt.printTree()
	}
}
