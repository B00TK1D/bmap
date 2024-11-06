# bmap
Better Map for Golang

Provides the following features over other standard library map implementations:

* Concurrency safe: all operations are protected using a RWMutex
* Ordered: map key-value pairs retain insertion order by default, which may be overriden by an arbitrary sort function
* Strongly typed: unlike `sync.Map`, bmaps enforce types (and don't require typecasting every return value)

Additionally, the `asyncBmap` provides ordered asynchronous write operations, providing exactly the same access guarantees as `bmap`, but with faster write calls.

## Usage

```go
package main

import (
	"fmt"
	"github.com/B00TK1D/bmap"
	"strings"
)

func main() {
	bmap := bmap.Type[string, int]{}.NewAsync()
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
```
