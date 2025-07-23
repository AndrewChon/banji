package banji

import (
	"cmp"
	"sync"
	"sync/atomic"
)

var queueIDCounter atomic.Uint64

// A Queue is a concurrency-safe min-priority queue that implements a pairing heap.
// in a concurrent environment, the locking order must be consistent. Lock queues in order of their IDs (Queue of id 1,
// ought to be locked before Queue of id 4, for example).
type Queue[K cmp.Ordered, V any] struct {
	sync.RWMutex
	id uint64

	root *tree[K, V]
	size int
}

func NewQueue[K cmp.Ordered, V any]() *Queue[K, V] {
	return &Queue[K, V]{
		id:   queueIDCounter.Add(1),
		root: nil,
		size: 0,
	}
}

type tree[K cmp.Ordered, V any] struct {
	key K
	val V

	parent           *tree[K, V]
	nextOlderSibling *tree[K, V]
	youngestChild    *tree[K, V]
}

func (q *Queue[K, V]) Push(elem V, priority K) {
	q.Lock()
	defer q.Unlock()

	q.root, _ = insert(q.root, priority, elem)
	q.size++
}

func (q *Queue[K, V]) Pop() (val V, ok bool) {
	q.Lock()
	defer q.Unlock()

	t := findMin(q.root)
	if t == nil {
		return
	}

	val = t.val
	q.root = removeMin(q.root)
	q.size--

	return val, true
}

func (q *Queue[K, V]) Peek() V {
	q.RLock()
	defer q.RUnlock()

	return findMin(q.root).val
}

func (q *Queue[K, V]) Size() int {
	q.RLock()
	defer q.RUnlock()

	return q.size
}

func (q *Queue[K, V]) Clear() {
	q.Lock()
	defer q.Unlock()

	q.root = nil
	q.size = 0
}

func (q *Queue[K, V]) Meld(other *Queue[K, V]) {
	if q.id < other.id {
		q.Lock()
		other.Lock()
	} else {
		other.Lock()
		q.Lock()
	}

	defer q.Unlock()
	defer other.Unlock()

	q.root = meld(q.root, other.root)
	q.size += other.size

	other.root = nil
	other.size = 0
}

// Pairing Heap implementation

// FindMin returns the root node of the tree, or nil if the tree is empty.
func findMin[K cmp.Ordered, V any](t *tree[K, V]) *tree[K, V] {
	if t == nil {
		return nil
	}
	return t
}

func (t *tree[K, V]) addChild(ct *tree[K, V]) {
	ct.parent = t

	if t.youngestChild == nil {
		t.youngestChild = ct
		ct.nextOlderSibling = nil
	} else {
		ct.nextOlderSibling = t.youngestChild
		t.youngestChild = ct
	}
}

// Meld forms a new tree from two other trees, with the largest becoming parent to the smallest.
func meld[K cmp.Ordered, V any](a, b *tree[K, V]) *tree[K, V] {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	if a.key < b.key {
		a.addChild(b)
		return a
	}
	b.addChild(a)
	return b
}

// Insert creates a new node with the specified key and value and melds it into the specified tree.
func insert[K cmp.Ordered, V any](h *tree[K, V], key K, val V) (newH, newNode *tree[K, V]) {
	newNode = &tree[K, V]{
		key: key,
		val: val,
	}
	newH = meld(h, newNode)

	return newH, newNode
}

// TwoPassMerge reconstitutes a rootless tree given its youngest child (and by extension, all of its children) into a
// new tree, with its smallest member as its Root.
func twoPassMerge[K cmp.Ordered, V any](yc *tree[K, V]) *tree[K, V] {
	if yc == nil || yc.nextOlderSibling == nil {
		return yc
	}

	var stack []*tree[K, V]

	// Meld the siblings in pairs, pairing the youngest sibling with the next older sibling.
	cur := yc
	for cur != nil {
		a := cur
		b := cur.nextOlderSibling

		if b != nil {
			cur = b.nextOlderSibling
			a.nextOlderSibling = nil
			b.nextOlderSibling = nil
			stack = append(stack, meld(a, b))
		} else {
			cur = nil
			a.nextOlderSibling = nil
			stack = append(stack, a)
		}
	}

	// Meld together the first-pass pairs, but in the opposite direction to prevent the overall tree from becoming
	// lopsided. The resulting tree will now have the smallest as its Root.
	for i := len(stack) - 2; i >= 0; i-- {
		stack[i] = meld(stack[i], stack[i+1])
	}

	return stack[0]
}

// RemoveMin removes the root node (in other words, the smallest node) from the provided tree, rebuilds the tree, and
// returns the new root node.
func removeMin[K cmp.Ordered, V any](t *tree[K, V]) *tree[K, V] {
	if t == nil {
		return nil
	}

	// We explicitly "disown" all its children. While not strictly necessary, it can free up the original Root node for
	// GC earlier, and proactively prevent any strange edge cases that may occur.
	youngestChild := t.youngestChild
	current := youngestChild
	for current != nil {
		next := current.nextOlderSibling
		current.parent = nil
		current = next
	}

	return twoPassMerge(youngestChild)
}

// Emancipate is a helper function that detaches a node from its parent.
func emancipate[K cmp.Ordered, V any](t *tree[K, V]) {
	parent := t.parent
	if parent.youngestChild == t {
		parent.youngestChild = t.nextOlderSibling
	} else {
		ys := parent.youngestChild
		for ys != nil && ys.nextOlderSibling != t {
			ys = ys.nextOlderSibling
		}
		if ys != nil {
			ys.nextOlderSibling = t.nextOlderSibling
		}
	}

	t.parent = nil
	t.nextOlderSibling = nil
}

// DecreaseKey decreases the target node's key to the provided new key. The new key must be less than the node's current
// key.
func decreaseKey[K cmp.Ordered, V any](h *tree[K, V], target *tree[K, V], newKey K) (newH *tree[K, V]) {
	if h == nil || target == nil {
		return h
	}
	if newKey >= target.key {
		return h
	}
	target.key = newKey

	if target.parent == nil {
		return h
	}
	emancipate(target)

	h = meld(h, target)
	return h
}
