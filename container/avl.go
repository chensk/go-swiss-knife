package container

import (
	"fmt"
	"math"
	"strings"
)

// BalancedBinarySearchTree implements balanced binary search tree (AVL tree).
type BalancedBinarySearchTree struct {
	root *avlTreeNode
	size int
}

// NewBinarySearchTree creates AVL tree from slice.
func NewBinarySearchTree(values []NodeValue) (*BalancedBinarySearchTree, error) {
	tree := &BalancedBinarySearchTree{}
	if len(values) == 0 {
		return tree, nil
	}

	for _, v := range values {
		tree.Insert(v)
	}
	return tree, nil
}

// Height returns height of tree.
func (t *BalancedBinarySearchTree) Height() int {
	_, h := _validateAvl(t.root)
	return h
}

// Size returns size of tree.
func (t *BalancedBinarySearchTree) Size() int {
	return t.size
}

// ValidateAvl returns whether the tree is AVL tree which can be asserted true.
func (t *BalancedBinarySearchTree) ValidateAvl() bool {
	b, _ := _validateAvl(t.root)
	return b
}

// Insert inserts a new node and keeps the tree an AVL tree. It should take O(logN) time.
func (t *BalancedBinarySearchTree) Insert(value NodeValue) {
	ok, _, r := avlInsert(t.root, t.root, value, nil, true)
	if ok {
		t.size++
		t.root = r
	}
}

// Delete deletes specified node and keeps the tree an AVL tree. It should take O(logN) time.
func (t *BalancedBinarySearchTree) Delete(value NodeValue) bool {
	ok, _, r := avlDelete(t.root, t.root, nil, value)
	if ok {
		t.size--
		t.root = r
	}
	return ok
}

// Traverse traverse tree with specified order and call the function for each non-nil node.
func (t *BalancedBinarySearchTree) Traverse(f func(value NodeValue), order TraverseOrder) {
	avlTraverse(t.root, f, order)
}

// String returns the inorder sequence and preorder sequence.
func (t *BalancedBinarySearchTree) String() string {
	inOrderResult := make([]NodeValue, 0, t.size)
	preOrderResult := make([]NodeValue, 0, t.size)
	t.Traverse(func(value NodeValue) {
		inOrderResult = append(inOrderResult, value)
	}, InOrder)
	t.Traverse(func(value NodeValue) {
		preOrderResult = append(preOrderResult, value)
	}, PreOrder)
	return fmt.Sprintf("inorder: %+v, preorder: %+v", inOrderResult, preOrderResult)
}

// PrettyPrint returns the formatted string of the tree.
func (t *BalancedBinarySearchTree) PrettyPrint() string {
	if t.root == nil {
		return "empty"
	}
	type Ele struct {
		tab int
		n   *avlTreeNode
	}
	stack := make([]Ele, 0)
	stack = append(stack, Ele{tab: 0, n: t.root})
	lines := make([]string, 0, t.size)
	for len(stack) > 0 {
		p := stack[0]
		stack = stack[1:]
		var item string = middleItem
		if len(stack) == 0 || stack[0].tab != p.tab {
			item = lastItem
		}
		var text = "nil"
		if p.n != nil {
			text = p.n.value.String()
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

// Exist returns whether the value exists.
func (t *BalancedBinarySearchTree) Exist(value NodeValue) bool {
	found, _, _ := avlSearch(value, t.root, nil)
	return found
}

func avlDelete(root *avlTreeNode, node *avlTreeNode, parent *avlTreeNode, value NodeValue) (bool, bool, *avlTreeNode) {
	if node == nil {
		return false, false, root
	}
	if node.value.Compare(value) < 0 {
		if ok, shorten, r := avlDelete(root, node.right, node, value); !ok {
			return false, false, r
		} else if shorten {
			switch node.bf {
			case LH:
				r, s := avlLeftBalance(root, node, parent)
				return true, s, r
			case EH:
				node.bf = LH
				return true, false, r
			case RH:
				node.bf = EH
				return true, true, r
			}
		} else {
			return true, false, r
		}
	} else if node.value.Compare(value) > 0 {
		if ok, shorten, r := avlDelete(root, node.left, node, value); !ok {
			return false, false, r
		} else if shorten {
			switch node.bf {
			case LH:
				node.bf = EH
				return true, true, r
			case EH:
				node.bf = RH
				return true, false, r
			case RH:
				r, s := avlRightBalance(root, node, parent)
				return true, s, r
			}
		} else {
			return true, false, r
		}
	}

	if node.left == nil {
		if parent == nil {
			return true, true, node.right
		} else if node == parent.left {
			parent.left = node.right
		} else {
			parent.right = node.right
		}
		return true, true, root
	} else if node.right == nil {
		if parent == nil {
			return true, true, node.left
		} else if node == parent.left {
			parent.left = node.left
		} else {
			parent.right = node.left
		}
		return true, true, root
	} else {
		var q = node
		var s = node.left
		shorten, r := _recursiveQ(q, s, node, root, parent)
		return true, shorten, r
	}
}

func _recursiveQ(q *avlTreeNode, s *avlTreeNode, node *avlTreeNode, root *avlTreeNode, parent *avlTreeNode) (bool, *avlTreeNode) {
	if s.right != nil {
		if shorten, r := _recursiveQ(s, s.right, node, root, q); shorten {
			if q == node {
				switch q.bf {
				case RH:
					rr, b := avlRightBalance(root, q, parent)
					return b, rr
				case EH:
					q.bf = RH
					return false, r
				case LH:
					q.bf = EH
					return true, r
				default:
					panic("unexpected")
				}
			}
			switch q.bf {
			case LH:
				rr, b := avlLeftBalance(root, q, parent)
				return b, rr
			case EH:
				q.bf = LH
				return false, r
			case RH:
				q.bf = EH
				return true, r
			default:
				panic("unexpected")
			}
		} else {
			return false, r
		}
	}
	node.value = s.value
	if q == node {
		q.left = s.left
		switch node.bf {
		case LH:
			node.bf = EH
			return true, root
		case EH:
			node.bf = RH
			return false, root
		case RH:
			node.bf = EH
			r, sh := avlRightBalance(root, node, parent)
			return sh, r
		default:
			panic("expected")
		}
	} else {
		q.right = s.left
		switch q.bf {
		case LH:
			r, sh := avlLeftBalance(root, q, parent)
			return sh, r
		case EH:
			q.bf = LH
			return false, root
		case RH:
			q.bf = EH
			return true, root
		default:
			panic("unexpected")
		}
	}
}

func avlInsert(root *avlTreeNode, node *avlTreeNode, value NodeValue, parent *avlTreeNode, left bool) (bool, bool, *avlTreeNode) {
	if root == nil {
		return true, true, &avlTreeNode{value: value}
	}
	if node == nil {
		ele := &avlTreeNode{value: value}
		if left {
			parent.left = ele
		} else {
			parent.right = ele
		}
		return true, true, root
	}
	if node.value.Compare(value) == 0 {
		return false, false, node
	}
	if node.value.Compare(value) > 0 {
		if ok, taller, _ := avlInsert(root, node.left, value, node, true); !ok {
			return false, false, root
		} else if taller {
			switch node.bf {
			case LH:
				r, _ := avlLeftBalance(root, node, parent)
				return true, false, r
			case EH:
				node.bf = LH
				return true, true, root
			case RH:
				node.bf = EH
				return true, false, root
			default:
				panic("unexpected")
			}
		} else {
			return true, false, root
		}
	} else {
		if ok, taller, _ := avlInsert(root, node.right, value, node, false); !ok {
			return false, false, root
		} else if taller {
			switch node.bf {
			case RH:
				r, _ := avlRightBalance(root, node, parent)
				return true, false, r
			case EH:
				node.bf = RH
				return true, true, root
			case LH:
				node.bf = EH
				return true, false, root
			default:
				panic("unexpected")
			}
		} else {
			return true, false, root
		}
	}
}

func avlSearch(value NodeValue, root *avlTreeNode, lastVisited *avlTreeNode) (bool, *avlTreeNode, *avlTreeNode) {
	if root == nil {
		return false, nil, lastVisited
	}
	if root.value.Compare(value) == 0 {
		return true, root, lastVisited
	}
	if root.value.Compare(value) > 0 {
		return avlSearch(value, root.left, root)
	}
	return avlSearch(value, root.right, root)
}

func avlTraverse(root *avlTreeNode, f func(value NodeValue), order TraverseOrder) {
	if root == nil {
		return
	}
	switch order {
	case PreOrder:
		f(root.value)
		avlTraverse(root.left, f, order)
		avlTraverse(root.right, f, order)
	case InOrder:
		avlTraverse(root.left, f, order)
		f(root.value)
		avlTraverse(root.right, f, order)
	case PostOrder:
		avlTraverse(root.left, f, order)
		avlTraverse(root.right, f, order)
		f(root.value)
	default:
		panic("unexpected order")
	}
}

func avlLeftBalance(root *avlTreeNode, node *avlTreeNode, parent *avlTreeNode) (*avlTreeNode, bool) {
	if node == nil {
		return root, false
	}
	// assert: current bf is larger than 1
	switch node.left.bf {
	case LH:
		node.bf = EH
		node.left.bf = EH
		return avlRightRotate(root, node, parent), true
	case RH:
		lr := node.left.right
		switch lr.bf {
		case LH:
			node.bf = RH
			node.left.bf = EH
		case RH:
			node.bf = EH
			node.left.bf = LH
		case EH:
			node.bf = EH
			node.left.bf = EH
		}
		lr.bf = EH
		root = avlLeftRotate(root, node.left, node)
		return avlRightRotate(root, node, parent), true
	case EH:
		node.bf = LH
		node.left.bf = RH
		return avlRightRotate(root, node, parent), false
	default:
		panic("unexpected")
	}
}

func avlRightBalance(root *avlTreeNode, node *avlTreeNode, parent *avlTreeNode) (*avlTreeNode, bool) {
	if node == nil {
		return root, false
	}
	// assert: current bf is larger than 1
	switch node.right.bf {
	case RH:
		node.bf = EH
		node.right.bf = EH
		return avlLeftRotate(root, node, parent), true
	case LH:
		rl := node.right.left
		switch rl.bf {
		case RH:
			node.bf = LH
			node.right.bf = EH
		case LH:
			node.bf = EH
			node.right.bf = RH
		case EH:
			node.bf = EH
			node.right.bf = EH
		}
		rl.bf = EH
		root = avlRightRotate(root, node.right, node)
		return avlLeftRotate(root, node, parent), true
	case EH:
		node.bf = RH
		node.right.bf = LH
		return avlLeftRotate(root, node, parent), false
	default:
		panic("unexpected")
	}
}

// left rotate tree left and return new root
func avlLeftRotate(root *avlTreeNode, node *avlTreeNode, parent *avlTreeNode) *avlTreeNode {
	r := node.right
	node.right = r.left
	r.left = node
	// case 0: parent is absent, meaning node is root
	if parent == nil {
		return r
	}
	// case 1: parent is present
	if parent.left == node {
		parent.left = r
	} else {
		parent.right = r
	}
	return root
}

// left rotate tree right and return new root
func avlRightRotate(root *avlTreeNode, node *avlTreeNode, parent *avlTreeNode) *avlTreeNode {
	l := node.left
	node.left = l.right
	l.right = node
	// case 0: parent is absent, meaning node is root
	if parent == nil {
		return l
	}
	// case 1: parent is present
	if parent.left == node {
		parent.left = l
	} else {
		parent.right = l
	}
	return root
}

func _validateAvl(root *avlTreeNode) (bool, int) {
	if root == nil {
		return true, 0
	}
	lb, lh := _validateAvl(root.left)
	rb, rh := _validateAvl(root.right)
	return lb && rb && math.Abs(float64(lh)-float64(rh)) <= 1, int(math.Max(float64(lh), float64(rh))) + 1
}

type avlTreeNode struct {
	value NodeValue
	left  *avlTreeNode
	right *avlTreeNode
	bf    int
}

// NodeValue represents the element type storing in tree.
type NodeValue interface {
	// Compare method compares with another NodeValue and return 1 if greater, 0 if equal, -1 if less.
	Compare(value NodeValue) int

	// NodeValue extends fmt.Stringer whose String() method would return formatted string.
	fmt.Stringer
}

type TraverseOrder int
type BalanceFactor int

const (
	PreOrder TraverseOrder = iota
	InOrder
	PostOrder

	LH = 1
	EH = 0
	RH = -1

	newLine    = "\n"
	emptySpace = "    "
	middleItem = "├── "
	lastItem   = "└── "
)
