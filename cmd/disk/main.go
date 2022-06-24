package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

// common
var (
	PageSize         = uintptr(os.Getpagesize())
	RowSize  uintptr = 100
)

type NodeType uint8

const (
	NodeInternal NodeType = iota + 1
	NodeLeaf
)

/*
 * Common Node Header Layout
 * 1. 节点类型，叶子节点、内部节点, uint8
 * 2. 是否根节点，uint8
 * 3. 父节点Id指针，uint32
 */
var (
	NodeTypeSize         = unsafe.Sizeof(uint8(0))
	NodeTypeOffset       = 0
	IsRootSize           = unsafe.Sizeof(uint8(0))
	IsRootOffset         = NodeTypeSize
	ParentPointerSize    = unsafe.Sizeof(uint32(0))
	ParentPointerOffset  = IsRootOffset + IsRootSize
	CommonNodeHeaderSize = NodeTypeSize + IsRootSize + ParentPointerSize
)

/*
 * Leaf Node Header Layout
 */
var (
	LeafNodeNumCellsSize   = unsafe.Sizeof(uint32(0))
	LeafNodeNumCellsOffset = CommonNodeHeaderSize
	LeafNodeHeaderSize     = CommonNodeHeaderSize + LeafNodeNumCellsSize
	LeafNodeNextLeafOffset = LeafNodeNumCellsOffset + LeafNodeNumCellsSize
)

/*
 * Leaf Node Body Layout
 */
var (
	LeafNodeKeySize               = unsafe.Sizeof(uint32(0))
	LeafNodeKeyOffset     uintptr = 0
	LeafNodeValueSize             = RowSize
	LeafNodeValueOffset           = LeafNodeKeyOffset + LeafNodeKeySize
	LeafNodeCellSize              = LeafNodeKeySize + LeafNodeValueSize
	LeafNodeSpaceForCells         = PageSize - LeafNodeHeaderSize
	LeafNodeMaxCells              = LeafNodeSpaceForCells / LeafNodeCellSize
)

/*
 * Internal Node Header Layout
 */
var (
	InternalNodeNumKeysSize      = unsafe.Sizeof(uint32(0))
	InternalNodeNumKeysOffset    = CommonNodeHeaderSize
	InternalNodeRightChildSize   = unsafe.Sizeof(uint32(0)) // 右孩子
	InternalNodeRightChildOffset = InternalNodeNumKeysOffset + InternalNodeNumKeysSize
	InternalNodeHeaderSize       = CommonNodeHeaderSize + InternalNodeNumKeysSize + InternalNodeRightChildSize
)

/*
 * Internal Node Body Layout
 */
var (
	InternalNodeKeySize   = unsafe.Sizeof(uint32(0))
	InternalNodeChildSize = unsafe.Sizeof(uint32(0))
	InternalNodeCellSize  = InternalNodeChildSize + InternalNodeKeySize
)

func getNodeType(node uintptr) NodeType {
	value := *(*uint8)(unsafe.Pointer(node + uintptr(NodeTypeOffset)))
	return NodeType(value)
}

func setNodeType(node uintptr, typ NodeType) {
	*(*uint8)(unsafe.Pointer(node + uintptr(NodeTypeOffset))) = uint8(typ)
}

func isNodeRoot(node uintptr) bool {
	value := *(*uint8)(unsafe.Pointer(node + IsRootOffset))
	return value == 1
}

func setNodeRoot(node uintptr, isRoot bool) {
	value := uint8(0)
	if isRoot {
		value = 1
	}
	*(*uint8)(unsafe.Pointer(node + IsRootOffset)) = value
}

func nodeParent(node uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(node + ParentPointerOffset))
}

func internalNodeNumKeys(node uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(node + InternalNodeNumKeysOffset))
}

func internalNodeRightChild(node uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(node + InternalNodeRightChildSize))
}

func internalNodeCell(node uintptr, cellNum uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(node + InternalNodeHeaderSize + uintptr(cellNum)*InternalNodeCellSize))
}

func internalNodeChild(node uintptr, childNum uint32) *uint32 {
	numKeys := *internalNodeNumKeys(node)
	if childNum > numKeys {
		log.Fatalf("Tried to access child_num %d > num_keys %d\n", childNum, numKeys)
		return nil
	} else if childNum == numKeys {
		return internalNodeRightChild(node)
	} else {
		return internalNodeCell(node, childNum)
	}
}

func internalNodeKey(node uintptr, keyNum uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(*internalNodeCell(node, keyNum)) + InternalNodeChildSize))
}

func leafNodeNumCells(node uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(node + LeafNodeNumCellsOffset))
}

func leafNodeCell(node uintptr, cellNum uint32) uintptr {
	return node + LeafNodeHeaderSize + uintptr(cellNum)*LeafNodeCellSize
}

func leafNodeKey(node uintptr, cellNum uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(leafNodeCell(node, cellNum)))
}

func leafNodeValue(node uintptr, cellNum uint32) uintptr {
	return leafNodeCell(node, cellNum) + LeafNodeKeySize
}

func leafNodeNextLeaf(node uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(node + LeafNodeNextLeafOffset))
}

func getNodeMaxKey(node uintptr) uint32 {
	switch getNodeType(node) {
	case NodeInternal:
		return *internalNodeKey(node, *internalNodeNumKeys(node)-1)
	case NodeLeaf:
		return *leafNodeKey(node, *leafNodeNumCells(node)-1)
	}
	panic("invalid node type")
}

func initializeLeafNode(node uintptr) {
	setNodeType(node, NodeLeaf)
	setNodeRoot(node, false)
	*leafNodeNumCells(node) = 0
	*leafNodeNextLeaf(node) = 0 // 0 represents no sibling
}

func initializeInternalNode(node uintptr) {
	setNodeType(node, NodeInternal)
	setNodeRoot(node, false)
	*internalNodeNumKeys(node) = 0
}

func main() {
	fmt.Println("hell disk")
}
