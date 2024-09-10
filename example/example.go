package main

import (
  "fmt"
  "github.com/B00TK1D/bmap"
)

func main() {
  bmap := bmap.Type[string, int]{}.NewAsync()
  bmap.Set("test2", 3)
  bmap.Set("test1", 1)
  bmap.Set("test3", 7)

  a, _ := bmap.Get("test1")
  b, _ := bmap.Get("test2")
  fmt.Println(a+b)

  for key, val := range bmap.Range() {
    fmt.Println(key, val)
  } 
}
