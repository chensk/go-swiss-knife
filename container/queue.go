package container

import (
	"fmt"
	"strings"
)

// LimitedBinarySearchQueue is a queue that supports binary search, or an Binary-Search tree that supports FIFO operations.
type LimitedBinarySearchQueue struct {
	root     *linkedListRbTreeNode
	size     int
	capacity int
	head     *linkedListRbTreeNode
	tail     *linkedListRbTreeNode
}

// NewLimitedQueueBinarySearch creates red-black tree from slice.
func NewLimitedQueueBinarySearch(values []NodeValue, capacity int) *LimitedBinarySearchQueue {
	tree := &LimitedBinarySearchQueue{
		root:     &linkedListRbTreeNode{},
		capacity: capacity,
	}
	if len(values) == 0 {
		return tree
	}

	for _, v := range values {
		tree.Insert(v)
	}
	return tree
}

type linkedListRbTreeNode struct {
	*rbTreeNode
	prev *linkedListRbTreeNode
	next *linkedListRbTreeNode
}

func (t *LimitedBinarySearchQueue) Height() int {
	return rbHeight(t.root.rbTreeNode)
}

func (t *LimitedBinarySearchQueue) Size() int {
	return t.size
}

func (t *LimitedBinarySearchQueue) Insert(value NodeValue) {
	if t.Size() >= t.capacity {
		head := t.head
		t.head = t.head.next

		r, _, b := rbDelete(t.root.rbTreeNode, t.root.rbTreeNode, head.value)
		if b {
			t.size--
			t.root = &linkedListRbTreeNode{
				rbTreeNode: r,
				// prev and next is useless for root, ignore it
			}
		}
	}
	ok, r, inserted := rbInsert(t.root.rbTreeNode, value)
	if ok {
		tail := t.tail
		nn := &linkedListRbTreeNode{
			rbTreeNode: inserted,
			prev:       tail,
		}
		if t.head == nil {
			t.head = nn
		}
		t.size++
		t.root = &linkedListRbTreeNode{
			rbTreeNode: r,
			// prev and next is useless for root, ignore it
		}
		t.root.color = BLACK
		// enqueue
		if tail != nil {
			tail.next = nn
		}
		t.tail = nn
	}
}

func (t *LimitedBinarySearchQueue) Delete(value NodeValue) bool {
	// delete is not permitted for limited binary search queue
	panic("Delete not implemented!")
}

func (t *LimitedBinarySearchQueue) Validate() bool {
	if t.root == nil {
		return true
	}
	if t.root.color == RED {
		return false
	}
	lb, li := validateRb(t.root.left)
	rb, ri := validateRb(t.root.right)
	return lb && rb && li == ri
}

func (t *LimitedBinarySearchQueue) Traverse(f func(value NodeValue) bool, order TraverseOrder) {
	rbTraverse(t.root.rbTreeNode, func(node *rbTreeNode) bool {
		return f(node.value)
	}, order)
}

func (t *LimitedBinarySearchQueue) String() string {
	inOrderResult := make([]NodeValue, 0, t.size)
	preOrderResult := make([]NodeValue, 0, t.size)
	t.Traverse(func(value NodeValue) bool {
		inOrderResult = append(inOrderResult, value)
		return true
	}, InOrder)
	t.Traverse(func(value NodeValue) bool {
		preOrderResult = append(preOrderResult, value)
		return true
	}, PreOrder)
	return fmt.Sprintf("inorder: %+v, preorder: %+v", inOrderResult, preOrderResult)
}

func (t *LimitedBinarySearchQueue) PrettyPrint() string {
	if t.root == nil {
		return "empty"
	}
	type Ele struct {
		tab int
		n   *rbTreeNode
	}
	stack := make([]Ele, 0)
	stack = append(stack, Ele{tab: 0, n: t.root.rbTreeNode})
	lines := make([]string, 0, t.size)
	for len(stack) > 0 {
		p := stack[0]
		stack = stack[1:]
		var item = middleItem
		if len(stack) == 0 || stack[0].tab != p.tab {
			item = lastItem
		}
		var text = fmt.Sprint(black, "nil")
		if p.n != nil {
			if p.n.color == RED {
				text = fmt.Sprint(red, p.n.value.String())
			} else {
				text = fmt.Sprint(black, p.n.value.String())
			}
		}
		lines = append(lines, fmt.Sprintf("%s%s%s", strings.Repeat(emptySpace, p.tab), item, text))
		if p.n != nil {
			if p.n.left != nil || p.n.right != nil {
				stack = append([]Ele{{tab: p.tab + 1, n: p.n.right}}, stack...)
				stack = append([]Ele{{tab: p.tab + 1, n: p.n.left}}, stack...)
			}
		}
	}
	return newLine + strings.Join(lines, newLine)
}

func (t *LimitedBinarySearchQueue) Exist(value NodeValue) bool {
	found, _, _ := rbSearch(value, t.root.rbTreeNode, nil)
	return found
}
