package main

import (
	"fmt"
	"strconv"

	"github.com/pedrogao/btrees/bptree"
)

func main() {
	tree := bptree.NewBPTree()
	for i := 1; i <= 1000; i++ {
		tree.Insert(i, strconv.Itoa(i))
	}

	for i := 1; i <= 1000; i++ {
		tree.Delete(i)
	}

	graph := tree.Graph()
	fmt.Printf(graph)

}
