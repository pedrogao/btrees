package main

type internalNode struct {
	typ    uint8  // 节点类型
	isRoot uint8  // 是否为根节点
	parent uint32 // 父节点id
}
