package bptree

import (
	"fmt"
	"sort"
	"unsafe"

	"github.com/pedrogao/btrees/common"
)

type kv struct {
	key   int
	value string
}

type kvs []kv

func (a kvs) Len() int           { return len(a) }
func (a kvs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a kvs) Less(i, j int) bool { return a[i].key < a[j].key }

// leafNode 第一个 key 不为空
// +----++----++----++----+
// | k1 || v2 || k2 || v2 |
// +----++----++----++----+
type leafNode struct {
	kvs        kvs           // 内部kv对
	max, count int           // kv对数量
	next       *leafNode     // 下一个叶子节点
	p          *internalNode // 父节点
}

func newLeafNode(max int) *leafNode {
	return &leafNode{
		kvs: make([]kv, max),
		max: max,
	}
}

// find the index of a key in the leaf node.
// If the key exists in the node, it returns the index and true.
// If the key does not exist in the node, it returns index to
// insert the key (the index of the smallest key in the node that larger
// than the given key) and false.
func (l *leafNode) find(key int) (int, bool) {
	c := func(i int) bool {
		return l.kvs[i].key >= key
	}
	// count 很重要，表示搜索的右边界
	i := sort.Search(l.count, c)

	if i < l.count && l.kvs[i].key == key {
		return i, true
	}

	return i, false
}

func (l *leafNode) insert(key int, value string) {
	i, ok := l.find(key)
	// 不支持 key 重复，发现有 key 直接替换即可
	if ok {
		l.kvs[i].value = value
		return
	}
	copy(l.kvs[i+1:], l.kvs[i:l.count])
	l.kvs[i].key = key
	l.kvs[i].value = value
	l.count++
}

func (l *leafNode) split() *leafNode {
	next := newLeafNode(l.max)

	mid := l.getMinSize()
	copy(next.kvs, l.kvs[mid:])

	next.count = l.max - mid
	next.next = l.next

	l.count = l.count - (l.max - mid)
	l.next = next

	return next
}

func (l *leafNode) full() bool { return l.count >= l.max }

func (l *leafNode) halfFull() bool { return l.count >= l.max/2 }

func (l *leafNode) parent() *internalNode { return l.p }

func (l *leafNode) setParent(p *internalNode) { l.p = p }

func (l *leafNode) getMaxSize() int {
	return l.max
}

func (l *leafNode) getMinSize() int {
	return l.max / 2
}

func (l *leafNode) isRoot() bool {
	return l.p == nil
}

func (l *leafNode) getSize() int {
	return l.count
}

func (l *leafNode) remove(key int) bool {
	idx, b := l.find(key)
	if !b {
		return false
	}
	common.RemoveAt(l.kvs, idx)
	l.count--
	return true
}

func (l *leafNode) moveLastToFrontOf(n node) {
	other, ok := n.(*leafNode)
	if !ok {
		return
	}
	l.count--
	removeItem := common.RemoveAt(l.kvs, l.count)
	other.resize(1)
	copy(other.kvs[1:], other.kvs)
	other.kvs[0] = removeItem
}

func (l *leafNode) moveAllTo(neighbor node) {
	other, ok := neighbor.(*leafNode)
	if !ok {
		return
	}
	copy(other.kvs[other.count:], l.kvs)
	l.kvs = make([]kv, l.max)
	other.resize(l.count)
	l.count = 0
}

func (l *leafNode) isLeaf() bool {
	return true
}

func (l *leafNode) valueAt(i int) any {
	return l.kvs[i].value
}

func (l *leafNode) resize(i int) {
	l.count += i
}

func (l *leafNode) moveFirstToEndOf(n node) {
	other, ok := n.(*leafNode)
	if !ok {
		return
	}
	l.count--
	removeItem := common.RemoveAt(l.kvs, 0)
	other.kvs[other.count] = removeItem
	other.resize(1)
}

func (l *leafNode) getFirstKey() int {
	return l.kvs[0].key
}

func (l *leafNode) id() string {
	id := uintptr(unsafe.Pointer(l))
	return fmt.Sprintf("%x", id)
}

func (l *leafNode) nextNode() node {
	return l.next
}
