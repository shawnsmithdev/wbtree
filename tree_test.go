package wbtree

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

const (
	bigRatTestSize  = 50000
	cmpUintTestSize = 250000
)

func testTree() *Tree[*big.Rat, *big.Rat] {
	var test *Tree[*big.Rat, *big.Rat]
	two := big.NewRat(2, 1)
	one := big.NewRat(1, 1)
	six := big.NewRat(6, 1)

	test, _ = test.Insert(two, two)
	test, _ = test.Insert(one, one)
	test, _ = test.Insert(six, six)
	return test
}

func TestTreeGet(t *testing.T) {
	test := testTree()
	two := test.GetNode(big.NewRat(2, 1))
	if two == nil {
		t.Fatal()
	}
}

func TestTreeReplace(t *testing.T) {
	test := testTree()
	two := big.NewRat(2, 1)
	var added bool
	test, added = test.Insert(two, two)
	if added {
		t.Fatal()
	}
	dos := test.GetNode(two)
	if dos == nil || dos.key.Cmp(two) != 0 {
		t.Fatal()
	}
}

func TestTreeInsert(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var test *Tree[*big.Rat, string]
	for i := 0; i < bigRatTestSize; i++ {
		rat := big.NewRat(rand.Int63(), rand.Int63())
		test, _ = test.Insert(rat, rat.FloatString(6))
	}

	keys := test.Keys()
	for i, val := range keys {
		if i == 0 {
			continue
		}
		if val.Cmp(keys[i-1]) < 0 {
			t.Fatal()
		}
	}
}

type cmpUint uint64

func (c cmpUint) Cmp(other cmpUint) int {
	if c == other {
		return 0
	}
	if c < other {
		return -1
	}
	return 1
}

func testTree2() *Tree[cmpUint, cmpUint] {
	var test *Tree[cmpUint, cmpUint]
	test, _ = test.Insert(2, 2)
	test, _ = test.Insert(1, 1)
	test, _ = test.Insert(6, 6)
	return test
}

func TestTree2Get(t *testing.T) {
	test := testTree2()
	two := test.GetNode(2)
	if two == nil {
		t.Fatal()
	}
}

func TestTree2Replace(t *testing.T) {
	test := testTree2()
	var added bool
	test, added = test.Insert(2, 2)
	if added {
		t.Fatal()
	}
	dos := test.GetNode(2)
	if dos == nil || dos.key.Cmp(2) != 0 {
		t.Fatal()
	}
}

func TestTree2Insert(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var test *Tree[cmpUint, cmpUint]
	for i := 0; i < cmpUintTestSize; i++ {
		newVal := cmpUint(rand.Uint64())
		test, _ = test.Insert(newVal, newVal)

		if test.GetNode(newVal) == nil {
			t.Fatal("cannot find value just inserted")
		}
	}

	keys := test.Keys()
	for i, key := range keys {
		if i == 0 {
			continue
		}
		if test.GetNode(key) == nil {
			t.Fatal("should be able to find key")
		}
		if key < keys[i-1] {
			t.Fatal("should be in order")
		}
	}
}

func TestTree2Remove(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var test *Tree[cmpUint, cmpUint]

	for i := 0; i < cmpUintTestSize; i++ {
		newVal := cmpUint(rand.Uint64())
		test, _ = test.Insert(newVal, newVal)
	}
	keys := test.Keys()
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	var removed bool
	for _, key := range keys {
		if test.GetNode(key) == nil {
			t.Fatal("should be able to find key:" + fmt.Sprint(key))
		}
		if test, removed = test.Remove(key); !removed {
			t.Fatal("should be able to delete key:" + fmt.Sprint(key))
		}
	}
	if test != nil {
		t.Fatal("should have nothing left")
	}
}
