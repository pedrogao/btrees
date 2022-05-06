package bptree

import (
	"bytes"
	"fmt"
	"strconv"
)

// BPTree b+ tree
type BPTree struct {
	root        node
	maxLeaf     int
	maxInternal int
}

type Option func(tree *BPTree)

func MaxLeaf(max int) Option {
	return func(tree *BPTree) {
		tree.maxLeaf = max
	}
}

func MaxInternal(max int) Option {
	return func(tree *BPTree) {
		tree.maxInternal = max
	}
}

func NewBPTree(options ...Option) *BPTree {
	b := &BPTree{}
	for _, option := range options {
		option(b)
	}
	if b.maxInternal <= 0 {
		b.maxInternal = MaxKC
	}
	if b.maxLeaf <= 0 {
		b.maxLeaf = MaxKV
	}
	return b
}

// First returns the first leafNode
func (t *BPTree) First() *leafNode {
	var (
		tmp   = t.root
		inter *internalNode
		ok    bool
	)

	// 如果是内部节点，则一直往下找
	for tmp != nil {
		inter, ok = tmp.(*internalNode)
		if !ok {
			// tmp 不是内部节点，那么是叶子节点，直接 break
			break
		}
		tmp = inter.kcs[0].child
	}

	return tmp.(*leafNode)
}

func (t *BPTree) Empty() bool {
	return t.root == nil
}

// Insert key->value
func (t *BPTree) Insert(key int, value string) {
	// 如果是空树，那么新建 root 节点
	if t.root == nil {
		t.startRoot(key, value)
		return
	}
	// 非空，插入至叶子节点
	t.insertIntoLeaf(key, value)
}

// Delete key
func (t *BPTree) Delete(key int) {
	if t.Empty() {
		return
	}
	leaf := t.findLeaf(key)
	if leaf == nil {
		return
	}
	ok := leaf.remove(key)
	if !ok {
		return
	}
	t.coalesceOrRedistribute(leaf)
}

// Search searches the key in B+ tree
// If the key exists, it returns the value of key and true
// If the key does not exist, it returns an empty string and false
func (t *BPTree) Search(key int) (string, bool) {
	if t.Empty() {
		return "", false
	}

	leaf := t.findLeaf(key)
	if leaf == nil {
		return "", false
	}

	idx, b := leaf.find(key)
	if !b {
		return "", false
	}

	return leaf.kvs[idx].value, true
}

func (t *BPTree) coalesceOrRedistribute(n node) {
	if n.isRoot() {
		t.adjustRoot(n)
		return
	}
	// 如果是半满状态，无需分裂、重组
	if n.halfFull() {
		return
	}
	parent := n.parent()
	idx := parent.valueIndex(n)
	if idx < 0 {
		panic("can't find child")
	}
	var sibling node
	if idx == 0 {
		sibling = parent.kcs[idx+1].child
	} else {
		sibling = parent.kcs[idx-1].child
	}
	// 重组
	if n.getSize()+sibling.getSize() >= n.getMaxSize() {
		t.redistribute(sibling, n, parent, idx)
		return
	}
	// 合并
	if idx == 0 {
		// n 在左边，sibling 在右边
		t.coalesce(n, sibling, parent)
	} else {
		// n 在右边
		t.coalesce(sibling, n, parent)
	}
}

func (t *BPTree) redistribute(neighbor, n node,
	parent *internalNode, index int) {
	// 将 neighbor 末尾移到 node 的最前面
	// 或者将 node 的开始项移到 neighbor 末尾
	if index == 0 {
		neighbor.moveFirstToEndOf(n)
		// 更新父节点指针
		parent.setKeyAt(1, neighbor)
	} else {
		neighbor.moveLastToFrontOf(n)
		parent.setKeyAt(index, neighbor)
	}
}

func (t *BPTree) coalesce(neighbor, n node, parent *internalNode) {
	// 合并以后可能还需要合并或者重组
	// n 所有项移动到 neighbor
	n.moveAllTo(neighbor)
	// 从 parent 中删除 node
	parent.remove(n)
	t.coalesceOrRedistribute(parent)
}

func (t *BPTree) adjustRoot(oldRoot node) {
	// 根节点还不是最后一个节点，仍然是内部节点，且有一个孩子节点
	if oldRoot.getSize() == 1 && !oldRoot.isLeaf() {
		t.root = oldRoot.valueAt(0).(node)
		t.root.setParent(nil)
	}
	// 只剩下根节点了，且已经没有子节点了
	if oldRoot.isLeaf() && oldRoot.getSize() == 0 {
		t.root = nil
	}
}

func (t *BPTree) insertIntoLeaf(key int, value string) {
	leaf := t.findLeaf(key)
	if leaf == nil {
		return
	}
	leaf.insert(key, value)
	// leaf 是否需要分裂
	if !leaf.full() {
		return
	}
	// 叶子节点分裂，并将 key 插入父节点
	newNode := leaf.split()
	t.insertIntoParent(leaf, newNode, newNode.kvs[0].key)
}

func (t *BPTree) insertIntoParent(old, new node, firstKey int) {
	if old.isRoot() {
		// 新建 root，并替换 root
		root := newInternalNode(nil, t.maxInternal)
		root.insert(0, old)
		root.insert(firstKey, new)
		t.root = root
		return
	}
	parent := old.parent()
	new.setParent(parent)
	parent.insert(firstKey, new)
	// 父节点无需分裂
	if !parent.full() {
		return
	}
	parentSibling, midKey := parent.split()
	// 父节点仍需分裂
	t.insertIntoParent(parent, parentSibling, midKey)
}

func (t *BPTree) findLeaf(key int) *leafNode {
	var (
		tmp   = t.root
		inter *internalNode
		ok    bool
	)

	// 如果是内部节点，则一直往下找
	for tmp != nil {
		inter, ok = tmp.(*internalNode)
		if !ok {
			// tmp 不是内部节点，那么是叶子节点，直接 break
			break
		}
		tmp = inter.lookup(key)
	}

	return tmp.(*leafNode)
}

func (t *BPTree) startRoot(key int, value string) {
	n := newLeafNode(nil, t.maxLeaf)
	n.insert(key, value)
	t.root = n
}

func (t *BPTree) printGraph() {
	fmt.Println("-----------------------------------")
	cur := t.root
	if cur == nil {
		fmt.Println("empty tree")
	}
	t.printNode(cur)
}

func (t *BPTree) Graph() string {
	cur := t.root
	out := bytes.NewBufferString("")
	out.WriteString("digraph G {\n")
	t.graph(cur, out)
	out.WriteString("}\n")
	return out.String()
}

func (t *BPTree) graph(cur node, out *bytes.Buffer) {
	leafPrefix, internalPrefix := "LEAF_", "INT_"
	switch n := cur.(type) {
	case *leafNode:
		out.WriteString(leafPrefix)
		out.WriteString(n.id())
		out.WriteString("[shape=plain color=green ")
		out.WriteString("label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
		out.WriteString("<TR><TD COLSPAN=\"")
		out.WriteString(strconv.Itoa(n.getSize()))
		out.WriteString("\">P=")
		out.WriteString(n.id())
		out.WriteString("</TD></TR>\n")
		out.WriteString("<TR><TD COLSPAN=\"")
		out.WriteString(strconv.Itoa(n.getSize()))
		out.WriteString("\">")
		out.WriteString("max_size=")
		out.WriteString(strconv.Itoa(n.getMaxSize()))
		out.WriteString(",min_size=")
		out.WriteString(strconv.Itoa(n.getMinSize()))
		out.WriteString("</TD></TR>\n")
		out.WriteString("<TR>")

		for i := 0; i < n.count; i++ {
			out.WriteString(fmt.Sprintf("<TD>%d</TD>\n", n.kvs[i].key))
		}
		out.WriteString("</TR>")
		out.WriteString("</TABLE>>];\n")
		if n.next != nil {
			out.WriteString(leafPrefix)
			out.WriteString(n.id())
			out.WriteString(" -> ")
			out.WriteString(leafPrefix)
			out.WriteString(n.next.id())
			out.WriteString(";\n")
			out.WriteString("{rank=same ")
			out.WriteString(leafPrefix)
			out.WriteString(n.id())
			out.WriteString(" ")
			out.WriteString(leafPrefix)
			out.WriteString(n.next.id())
			out.WriteString("};\n")
		}
		if n.parent() != nil {
			out.WriteString(internalPrefix)
			out.WriteString(n.parent().id())
			out.WriteString(":p")
			out.WriteString(n.id())
			out.WriteString(" -> ")
			out.WriteString(leafPrefix)
			out.WriteString(n.id())
			out.WriteString(";\n")
		}
	case *internalNode:
		out.WriteString(internalPrefix)
		out.WriteString(n.id())
		out.WriteString("[shape=plain color=pink ")
		out.WriteString("label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
		out.WriteString("<TR><TD COLSPAN=\"")
		out.WriteString(strconv.Itoa(n.getSize()))
		out.WriteString("\">P=")
		out.WriteString(n.id())
		out.WriteString("</TD></TR>\n")
		out.WriteString("<TR><TD COLSPAN=\"")
		out.WriteString(strconv.Itoa(n.getSize()))
		out.WriteString("\">")
		out.WriteString("max_size=")
		out.WriteString(strconv.Itoa(n.getMaxSize()))
		out.WriteString(",min_size=")
		out.WriteString(strconv.Itoa(n.getMaxSize() / 2))
		out.WriteString("</TD></TR>\n")
		out.WriteString("<TR>")

		for i := 0; i < n.count; i++ {
			out.WriteString("<TD PORT=\"p")
			out.WriteString(n.kcs[i].child.id())
			out.WriteString("\">")
			if i > 0 {
				out.WriteString(strconv.Itoa(n.kcs[i].key))
			} else {
				out.WriteString(" ")
			}
			out.WriteString("</TD>\n")
		}
		out.WriteString("</TR>")
		out.WriteString("</TABLE>>];\n")
		if n.parent() != nil {
			out.WriteString(internalPrefix)
			out.WriteString(n.parent().id())
			out.WriteString(":p")
			out.WriteString(n.id())
			out.WriteString(" -> ")
			out.WriteString(internalPrefix)
			out.WriteString(n.id())
			out.WriteString(";\n")
		}

		for i := 0; i < n.count; i++ {
			t.graph(n.kcs[i].child, out)
			if i > 0 {
				isLeaf := n.kcs[i].child.isLeaf()
				if !isLeaf {
					out.WriteString("{rank=same ")
					out.WriteString(internalPrefix)
					out.WriteString(n.kcs[i-1].child.id())
					out.WriteString(" ")
					out.WriteString(internalPrefix)
					out.WriteString(n.kcs[i].child.id())
					out.WriteString("};\n")
				}
			}
		}
	}
}

func (t *BPTree) printTree() {
	fmt.Println("-----------------------------------")
	cur := t.root
	if cur == nil {
		fmt.Println("empty tree")
	}
	t.printNode(cur)
}

func (t *BPTree) printNode(cur node) {
	switch n := cur.(type) {
	case *leafNode:
		fmt.Printf("- leaf %s (size %d)\n", n.id(), n.count)
		for i := 0; i < n.count; i++ {
			fmt.Printf("<%d, %s>,", n.kvs[i].key, n.kvs[i].value)
		}
		fmt.Println()
		fmt.Println()
		break
	case *internalNode:
		fmt.Printf("- internal %s (size %d)\n", n.id(), n.count)
		for i := 0; i < n.count; i++ {
			key := n.kcs[i].key
			if i == 0 {
				key = 0
			}
			fmt.Printf("<%d, %s>", key, n.kcs[i].child.id())
		}
		fmt.Println()
		fmt.Println()
		for i := 0; i < n.count; i++ {
			child := n.kcs[i].child
			t.printNode(child)
		}
		break
	}
}
