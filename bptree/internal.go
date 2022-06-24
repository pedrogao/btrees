package bptree

import (
	"fmt"
	"sort"
	"unsafe"

	"github.com/pedrogao/btrees/common"
)

type kc struct {
	key   int
	child node
}

// one empty slot for split
type kcs []kc

func (a kcs) Len() int { return len(a) }

func (a kcs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a kcs) Less(i, j int) bool {
	if a[i].key == 0 {
		return false
	}

	if a[j].key == 0 {
		return true
	}

	return a[i].key < a[j].key
}

// internalNode 第一个 key 为空
// +---++----++----++----+
// |   || v1 || k1 || v2 |
// +---++----++----++----+
type internalNode struct {
	kcs        kcs           // kv键值对
	max, count int           // kv最大数量、数量
	p          *internalNode // 父节点
}

func newInternalNode(max int) *internalNode {
	// 但判断内部节点的标注仍以 key 为准，即以 key 作为 full，split 的标准
	i := &internalNode{
		max:   max,
		count: 0,
		kcs:   make([]kc, max),
	}

	return i
}

func (n *internalNode) find(key int) (int, bool) {
	// 如果没有任何数据，直接返回 0，false，即 0 号位插入，且未找到
	if n.count == 0 {
		return 0, false
	}
	// todo >= or >
	c := func(i int) bool { return n.kcs[i].key >= key }
	i := sort.Search(n.count, c)
	if i < n.count && n.kcs[i].key == key {
		return i, true
	}

	return i, false
}

func (n *internalNode) lookup(key int) node {
	// 如果没有任何数据，直接返回 0，false，即 0 号位插入，且未找到
	if n.count == 0 {
		return nil
	}

	c := func(i int) bool { return n.kcs[i].key >= key }
	i := sort.Search(n.count, c)

	if i < n.count && n.kcs[i].key == key {
		return n.kcs[i].child
	}

	if i == 0 {
		return n.kcs[0].child
	}

	if i < n.count && n.kcs[i].key > key {
		return n.kcs[i-1].child
	}

	if i >= n.count {
		return n.kcs[n.count-1].child
	}

	return nil
}

func (n *internalNode) full() bool { return n.count >= n.max }

func (n *internalNode) halfFull() bool { return n.count >= n.max/2 }

func (n *internalNode) parent() *internalNode { return n.p }

func (n *internalNode) setParent(p *internalNode) { n.p = p }

func (n *internalNode) insert(key int, child node) bool {
	// 即使 key 重复，仍然需要插入，因此 b+tree 的内部节点就是会重复的
	i, _ := n.find(key)
	if i >= n.max {
		return false
	}
	copy(n.kcs[i+1:], n.kcs[i:n.count])
	n.kcs[i].key = key
	n.kcs[i].child = child
	child.setParent(n)
	n.count++
	return true
}

func (n *internalNode) split() (*internalNode, int) {
	// 3/2 => 1
	midIndex := n.count / 2
	midKey := n.kcs[midIndex].key

	// create the split node without a parent
	next := newInternalNode(n.max)
	copy(next.kcs, n.kcs[midIndex:])
	next.count = n.count - midIndex
	// update parent
	for i := 0; i < next.count; i++ {
		next.kcs[i].child.setParent(next)
	}
	n.count = midIndex

	return next, midKey
}

func (n *internalNode) getMaxSize() int {
	return n.max
}

func (n *internalNode) getMinSize() int {
	return n.max / 2
}

func (n *internalNode) isRoot() bool {
	return n.p == nil
}

func (n *internalNode) getSize() int {
	return n.count
}

func (n *internalNode) valueIndex(val node) int {
	for i, item := range n.kcs {
		if item.child == val {
			return i
		}
	}
	return -1
}

func (n *internalNode) setKeyAt(index int, val node) {
	firstKey := val.getFirstKey()
	n.kcs[index].key = firstKey
	n.kcs[index].child = val
}

func (n *internalNode) moveLastToFrontOf(n2 node) {
	other, ok := n2.(*internalNode)
	if !ok {
		return
	}
	n.count--
	removeItem := common.RemoveAt(n.kcs, n.count)
	other.resize(1)
	copy(other.kcs[1:], other.kcs)
	other.kcs[0] = removeItem
}

func (n *internalNode) remove(n2 node) {
	idx := n.valueIndex(n2)
	common.RemoveAt(n.kcs, idx)
	n.count--
}

func (n *internalNode) moveAllTo(neighbor node) {
	other, ok := neighbor.(*internalNode)
	if !ok {
		return
	}
	copy(other.kcs[other.count:], n.kcs)
	n.kcs = make([]kc, n.max)
	other.resize(n.count)
	n.count = 0
}

func (n *internalNode) isLeaf() bool {
	return false
}

func (n *internalNode) valueAt(i int) any {
	return n.kcs[i].child
}

func (n *internalNode) resize(i int) {
	n.count += i
}

func (n *internalNode) moveFirstToEndOf(n2 node) {
	other, ok := n2.(*internalNode)
	if !ok {
		return
	}
	n.count--
	removeItem := common.RemoveAt(n.kcs, 0)
	other.kcs[other.count] = removeItem
	other.resize(1)
}

func (n *internalNode) getFirstKey() int {
	return n.kcs[0].key
}

func (n *internalNode) id() string {
	id := uintptr(unsafe.Pointer(n))
	return fmt.Sprintf("%x", id)
}

func (n *internalNode) nextNode() node {
	return nil
}
