package wbtree

// weight-balanced binary search tree
// see: https://yoichihirai.com/bst.pdf

const (
	// these params tune rebalance behavior, values found valid and performant in Y. Hirai (2011).
	treeDelta = 3
	treeGamma = 2
)

// Cmpable are types that can do C style comparisons with a value of the parameter type.
// This is almost always the same type.
// For example, *big.Rat and *big.Int implement this interface.
type Cmpable[T any] interface {
	// Cmp compares this value with other and returns:
	//   -1 if this value < other
	//    0 if this value == other
	//   +1 if this value > other
	Cmp(other T) int
}

// Tree is a map implemented with a weight-balanced tree.
type Tree[K Cmpable[K], V any] struct {
	left      *Tree[K, V] // lesser
	right     *Tree[K, V] // greater
	childSize uint64      // count of children, 0 for leaf nodes, equals size-1 and weight-2
	key       K
	value     V
}

// RootKey returns the key of the root node of this tree.
func (t *Tree[K, V]) RootKey() (key K) {
	if t != nil {
		key = t.key
	}
	return key
}

// RootValue returns the stored value in the root node of this tree.
func (t *Tree[K, V]) RootValue() (value V) {
	if t != nil {
		value = t.value
	}
	return value
}

// Size returns the total count of nodes in this tree, including the root node.
func (t *Tree[K, V]) Size() uint64 {
	if t == nil {
		return 0
	}
	return t.childSize + 1
}

// Keys returns a slice of all keys in dfs order
func (t *Tree[K, V]) Keys() []K {
	var keys []K
	if t.left != nil {
		keys = t.left.Keys()
	}
	keys = append(keys, t.key)
	if t.right != nil {
		rKeys := t.right.Keys()
		keys = append(keys, rKeys...)
	}
	return keys
}

// Values returns a slice of all values in dfs order
func (t *Tree[K, V]) Values() []V {
	var vals []V
	if t.left != nil {
		vals = t.left.Values()
	}
	vals = append(vals, t.value)
	if t.right != nil {
		rKeys := t.right.Values()
		vals = append(vals, rKeys...)
	}
	return vals
}

// GreatestNode returns the rightmost Tree node
func (t *Tree[K, V]) GreatestNode() *Tree[K, V] {
	if t == nil {
		return nil
	}
	if t.right == nil {
		return t
	}
	return t.right.GreatestNode()
}

// LeastNode returns the leftmost Tree node
func (t *Tree[K, V]) LeastNode() *Tree[K, V] {
	if t == nil {
		return nil
	}
	if t.left == nil {
		return t
	}
	return t.left.LeastNode()
}

// Get returns the value in this tree associated with key, or the zero value of V if key is not present
func (t *Tree[K, V]) Get(key K) V {
	return t.GetNode(key).RootValue()
}

// GetNode returns the Tree node at key, or nil if key is not present
func (t *Tree[K, V]) GetNode(key K) *Tree[K, V] {
	if t == nil {
		return nil
	}
	compared := t.key.Cmp(key)
	if compared == 0 {
		return t
	}
	var next *Tree[K, V]
	if compared < 0 {
		next = t.right
	} else {
		next = t.left
	}
	if next == nil {
		return nil
	}
	return next.GetNode(key)
}

// Insert adds a new value to this Tree, or updates an existing value.
// Returns the new root (the same node if not rebalanced).
// The returned bool is true if the Insert resulted in a new entry in the tree, false if an existing value was replaced.
func (t *Tree[K, V]) Insert(key K, val V) (*Tree[K, V], bool) {
	// if tree is empty, return new root
	if t == nil {
		return &Tree[K, V]{
			key:   key,
			value: val,
		}, true
	}

	// if key is equal to root, replace root value
	compared := t.key.Cmp(key)
	if compared == 0 {
		t.value = val
		return t, false
	}

	// add to right child if root key less than inserted key
	addRight := compared < 0
	var next *Tree[K, V]
	if addRight {
		next = t.right
	} else {
		next = t.left
	}

	// if no child, add new leaf node
	if next == nil {
		next = &Tree[K, V]{
			key:   key,
			value: val,
		}
	} else {
		// attempt to add to child tree
		var added bool
		next, added = next.Insert(key, val)
		if !added { // replaced in child tree
			return t, false
		}
	}

	// update child pointers and size
	if addRight {
		t.right = next
	} else {
		t.left = next
	}
	t.childSize++

	// fix balance and return
	return t.balance(addRight), true
}

// Remove removes a node from a tree with the given key, if any exists.
// Returns the new root node, and a bool that is true if a node was removed.
func (t *Tree[K, V]) Remove(key K) (*Tree[K, V], bool) {
	if t == nil {
		// cannot remove from the empty tree
		return t, false
	}
	compared := t.key.Cmp(key)
	if compared == 0 {
		// remove this node
		if t.left == nil && t.right == nil {
			// leaf node, return the empty tree
			return nil, true
		}

		if t.left != nil && t.right != nil {
			// two children, swap values with least on right
			rightMin := t.right.LeastNode()
			t.key = rightMin.key
			t.value = rightMin.value
			t.right, _ = t.right.Remove(rightMin.key)
			t.childSize--
			result := t.balance(false)
			return result, true
		}
		if t.left == nil {
			return t.right, true
		}
		return t.left, true
	}

	// not equal to t, determine next node to remove from
	removeRight := compared < 0
	var next *Tree[K, V]
	if removeRight {
		next = t.right
	} else {
		next = t.left
	}

	if newNext, removed := next.Remove(key); removed {
		if removeRight {
			t.right = newNext
		} else {
			t.left = newNext
		}
		t.childSize--
		result := t.balance(!removeRight)
		return result, true
	}

	// not present
	return t, false
}

// balance checks if, after Insert or Remove, the tree has become unbalanced.
// If it has, balance rotates (or double rotates) to fix it, updating childSize as needed.
// rightHeavy should be true if a new node was inserted into right (or removed from left), else false.
// Returns new root (same root if no rebalance is needed).
func (t *Tree[K, V]) balance(rightHeavy bool) *Tree[K, V] {
	x, c := t.left, t.right
	if !rightHeavy {
		x, c = c, x
	}
	var b, z *Tree[K, V]
	if c != nil {
		b, z = c.left, c.right
		if !rightHeavy {
			b, z = z, b
		}
	}

	xSize, cSize := x.Size(), c.Size()
	if (xSize+1)*treeDelta >= cSize+1 {
		// balanced, no rotation
		return t
	}

	// balance broken, check for single or double rotation
	bSize, zSize := b.Size(), z.Size()

	// single rotation
	if bSize+1 < (zSize+1)*treeGamma {
		t.left, t.right = x, b
		c.left, c.right = t, z
		if !rightHeavy {
			t.left, t.right = t.right, t.left
			c.left, c.right = c.right, c.left
		}
		t.childSize = xSize + bSize
		c.childSize = t.Size() + zSize
		return c
	}

	// double rotation
	s, y := b.left, b.right
	if !rightHeavy {
		s, y = y, s
	}

	t.left, t.right = x, s
	b.left, b.right = t, c
	c.left, c.right = y, z
	if !rightHeavy {
		t.left, t.right = t.right, t.left
		b.left, b.right = b.right, b.left
		c.left, c.right = c.right, c.left
	}

	t.childSize = xSize + s.Size()
	c.childSize = y.Size() + zSize
	b.childSize = t.Size() + c.Size()
	return b
}
