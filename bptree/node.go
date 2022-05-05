package bptree

const (
	MaxKV = 255
	MaxKC = 511
)

// TODO disk
type node interface {
	find(key int) (int, bool)
	parent() *internalNode
	setParent(*internalNode)
	full() bool
	halfFull() bool
	getMax() int
	getSize() int
	resize(int)
	isRoot() bool
	moveLastToFrontOf(node)
	moveAllTo(neighbor node)
	isLeaf() bool
	valueAt(i int) any
	moveFirstToEndOf(n node)
	getFirstKey() int
	id() string
	nextNode() node
}
