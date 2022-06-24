package b2

import (
	"github.com/pedrogao/btrees/common"
)

type item struct {
	key   int
	value any // value or child
}

func newItem(key int, val any) *item {
	return &item{
		key:   key,
		value: val,
	}
}

type node struct {
	parent   *node // 父节点
	max      int
	count    int
	items    []*item // 节点kv对
	children []*node // 子节点
}

func newNode(size int) *node {
	return &node{
		parent:   nil,
		max:      size,
		items:    make([]*item, size),
		children: make([]*node, size),
	}
}

type BTree struct {
	min  int
	max  int
	n    int
	root *node
}

func NewBTree(min int) *BTree {
	// min |= 1
	return &BTree{
		min: min,
		max: min * 2,
		n:   0,
	}
}

func (t *BTree) Insert(key int, val any) bool {
	cur := t.root
	if cur == nil {
		cur = newNode(t.max)
		t.root = cur
	}

	w := cur.add(key, val)
	if w != nil {
		n := newNode(t.max)
		last := cur.removeLast()
		n.addChild(0, nil, cur)
		n.addChild(last.key, last.value, w)
		t.root = n
	}

	t.n += 1
	return true
}

func (t *BTree) Search(key int) any {
	u := t.root
	var z any
	for u != nil {
		i := u.findIndex(key)
		if i < 0 { // found
			return u.items[-(i + 1)].value
		}
		if u.items[i] != nil {
			z = u.items[i].value
		}
		// search at sub node
		if i >= u.count {
			i = u.count - 1
		}
		u = u.children[i]
	}
	return z
}

func (t *BTree) Delete(key int) bool {
	r := t.root
	if r.delete(key) {
		t.n--
		if r.getSize() == 0 && t.n > 0 {
			t.root = t.root.children[0]
			t.root.parent = nil
		}
		return true
	}
	return false
}

func (n *node) delete(key int) bool {
	i := n.findIndex(key)
	if i < 0 {
		// found
		i = -(i + 1)
		if n.isLeaf() {
			n.remove(i)
		} else {
			smallest := n.children[i].removeSmallest()
			n.items[i] = smallest
			n.checkUnderflow(i)
		}
		return true
	}
	if i > 0 {
		i -= 1
	}
	// 从子节点中删除
	if n.children[i].delete(key) {
		// 判断是否需要重组、合并
		n.checkUnderflow(i)
		return true
	}

	return false
}

func (n *node) removeSmallest() *item {
	if n.isLeaf() {
		return n.remove(0)
	}
	y := n.children[0].removeSmallest()
	n.checkUnderflow(0)
	return y
}

func (n *node) checkUnderflow(i int) {
	if n.children[i] == nil {
		return
	}
	if i == 0 {
		// 如果被删除的节点是第一个节点，因此其 sibling 是右边的兄弟
		n.checkUnderflowZero(i)
	} else {
		// 如果删除的节点不是第一个节点，其 sibling 是左边的兄弟
		n.checkUnderflowNonZero(i)
	}
}

func (n *node) checkUnderflowZero(i int) {
	w := n.children[i]
	if w.underflow() {
		v := n.children[i+1]
		if v.getSize() > n.max/2 {
			leftRotation(n, v, w, i)
		} else {
			merge(n, v, w, i)
			n.children[i] = w
		}
	}
}

func (n *node) checkUnderflowNonZero(i int) {
	w := n.children[i]
	if w.underflow() {
		v := n.children[i-1] // sibling
		if v.getSize() > n.max/2 {
			// 如果 sibling 半满，那么可以借一个，否则只能合并
			rightRotation(n, v, w, i-1) // fixme
		} else {
			// 合并
			merge(n, v, w, i)
		}
	}
}

func merge(parent, v, w *node, i int) {
	sv := v.getSize()
	sw := w.getSize()
	// 合并孩子节点
	copy(v.items[sv+1:], w.items[:sw])
	v.count += sw
	copy(v.children[sv+1:], w.children[:sw+1]) // fixme
	v.items[sv] = parent.items[i]
	v.count += 1
	// 处理 parent
	if i == 0 {
		common.RemoveAt(parent.items, 0)
		common.RemoveAt(parent.items, 0) // 移除一个之后，i-1也变成了0
		parent.count -= 2
	} else {
		common.RemoveAt(parent.items, i-1)
		common.RemoveAt(parent.items, i-1)
		parent.count -= 2
	}
}

// leftRotation 左旋，右孩子节点减少一个项，替换父指针，然后父指针项补充到左孩子
// 然后重新达到平衡
func leftRotation(parent, v, w *node, i int) {
	sv := v.getSize()
	sw := w.getSize()
	shift := ((sw + sv) / 2) - sw
	w.items[shift] = parent.items[i+1]
	w.count++
	smallest := v.remove(0)
	parent.items[i+1] = smallest
}

func rightRotation(parent, v, w *node, i int) {
	sv := v.getSize()
	sw := w.getSize()
	shift := ((sw + sv) / 2) - sw
	w.items[sw] = parent.items[i]
	copy(w.items[sw+1:shift+sw], v.items)
	copy(w.children[sw+1:shift+sw+1], v.children)
	parent.items[i] = v.items[shift-1]
	copy(v.items, v.items[shift:v.max])
	copy(v.children, v.children[shift:v.max+1]) // ?fixme
}

func (n *node) findIndex(key int) int {
	lo, hi := 0, n.count
	for hi != lo {
		m := (hi + lo) / 2
		if n.items[m] == nil || key < n.items[m].key {
			hi = m
		} else if key > n.items[m].key {
			lo = m + 1
		} else {
			return -m - 1
		}
	}
	return lo
}

func (n *node) add(key int, val any) *node {
	return n.addChild(key, val, nil)
}

func (n *node) addChild(key int, val any, child *node) *node {
	i := n.findIndex(key)
	if i < 0 {
		panic("duplicate key")
	}
	if n.isLeaf() || child != nil { // 叶子节点，直接加入即可
		n.items[i] = newItem(key, val)
		if child != nil {
			n.children[i] = child
			child.parent = n
		}
		n.count++
	} else {
		if i >= n.count {
			i = n.count - 1
		}
		// 新的子节点
		w := n.children[i].addChild(key, val, child)
		if w != nil {
			x := n.children[i].removeLast()
			n.addChild(x.key, x.value, w)
		}
	}

	if n.full() {
		return n.split()
	}

	return nil
}

func (n *node) split() *node {
	// 5/2 = 2
	m := n.max / 2
	other := newNode(n.max)
	other.parent = n.parent
	copy(other.items, n.items[m:])
	copy(other.children, n.children[m:])
	other.count = n.count - m
	n.count = m
	return other
}

func (n *node) remove(idx int) *item {
	n.count--
	common.RemoveAt(n.children, idx)
	return common.RemoveAt(n.items, idx)
}

func (n *node) removeLast() *item {
	return n.remove(n.count - 1)
}

func (n *node) getMax() int {
	return n.max
}

func (n *node) getSize() int {
	return n.count
}

func (n *node) full() bool {
	return n.count >= n.max
}

func (n *node) underflow() bool {
	return n.count < (n.max/2 - 1)
}

func (n *node) isLeaf() bool {
	return n.children[0] == nil
}
