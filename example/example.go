package main

import (
	"fmt"
	"github.com/B00TK1D/bmap"
	"strings"
)

func main() {
	bmap := Bmap[string, int]{}
	bmap.Set("test2", 3)
	bmap.Set("test1", 9)
	bmap.Set("test3", 7)

	fmt.Println("Strongly typed")
	a, _ := bmap.Get("test1")
	b, _ := bmap.Get("test2")
	fmt.Println("Sum:", a+b)
	fmt.Println()

	fmt.Println("Iterator")
	for key, val := range bmap.Range() {
		fmt.Println(key, val)
	}
	fmt.Println()

	fmt.Println("Key sort")
	bmap.SortKeys(func(k1, k2 string) bool {
		return strings.Compare(k1, k2) < 0
	})
	fmt.Println(bmap)

	fmt.Println("Value sort")
	bmap.Sort(func(v1, v2 int) bool {
		return v1 < v2
	})
	fmt.Println(bmap)

	fmt.Println("Key deletion")
	bmap.Delete("test2")
	fmt.Println(bmap)

	fmt.Println("Swap example")
	bmap.Swap("test2", "test3")
	fmt.Println(bmap)

	fmt.Println("Missing key deletion")
	bmap.Delete("missing")
	fmt.Println(bmap)

	fmt.Println("Length")
	fmt.Println(bmap.Len())
}
