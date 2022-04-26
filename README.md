[![Go Reference](https://pkg.go.dev/badge/github.com/shawnsmithdev/wbtree.svg)](https://pkg.go.dev/github.com/shawnsmithdev/wbtree) [![Go Report Card](https://goreportcard.com/badge/github.com/shawnsmithdev/wbtree)](https://goreportcard.com/report/github.com/shawnsmithdev/wbtree) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# wbtree
wbtree is a weight balanced binary search tree for go 1.18+

For more on weight balanced trees, see: https://yoichihirai.com/bst.pdf

Concurrent access
=================
No. This library does not support concurrent access (it is not "thread-safe").
This may be addressed in a future release. It may not. It is probably at least feasible to make it lock-free...

Balance parameters
==================
The choice of `<3,2>` as balance parameters here is mostly for the convienience of using simple integer values.
There's a somewhat faster setting, `<1+sqrt(2), sqrt(2)>`, which is not even rational.
The performance is quite close even with the integer params, so they are used, but it should be noted
that I've not benchmarked or even attempted any others yet.

Basic usage
===========

```go
package main

import (
	"fmt"
	"github.com/shawnsmithdev/wbtree"
	"math/big"
)

func main() {
	var example *wbtree.Tree[*big.Int, string]
	var inserted, removed bool

	// insert and update
	example, inserted = example.Insert(big.NewInt(5), "fie")
	fmt.Println(inserted) // true
	example, inserted = example.Insert(big.NewInt(5), "five")
	fmt.Println(inserted) // false
	example, _ = example.Insert(big.NewInt(4), "four")
	example, _ = example.Insert(big.NewInt(3), "three")

	// remove
	fmt.Println(example.Keys()) // 5, 4, 3
	fmt.Println(example.Values()) // []string{"three", "four", "five"}
	example, removed = example.Remove(big.NewInt(4))
	fmt.Println(removed) // true
	example, removed = example.Remove(big.NewInt(42))
	fmt.Println(false) // true
	fmt.Println(example.Values()) // []string{"three", "five"}
}
```
