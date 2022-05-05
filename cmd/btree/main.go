package main

import (
	"fmt"

	"github.com/pedrogao/btrees/btree"
)

func main() {
	minimumItemsInNode := btree.DefaultMinItems
	tree := btree.NewTree(minimumItemsInNode)
	value := "0"
	tree.Put(value, value)

	retVal := tree.Find(value)
	fmt.Printf("Returned value is :%v \n", retVal)

	tree.Remove(value)

	retVal = tree.Find(value)
	fmt.Print("Returned value is nil")
}
