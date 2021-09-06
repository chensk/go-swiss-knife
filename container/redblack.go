package container

import (
	"fmt"
	"math"
	"strings"
)

// RedBlackTree implements red-black tree.
// Refers to https://mp.weixin.qq.com/s/4sCnvWmW7-fOIlpNeIIjIw
type RedBlackTree struct {
	root *rbTreeNode
	size int
}

type rbTreeNode struct {
	value  NodeValue
	parent *rbTreeNode
	left   *rbTreeNode
	right  *rbTreeNode
	color  nodeColor
}

// NewRedBlackTree creates red-black tree from slice.
func NewRedBlackTree(values []NodeValue) (*RedBlackTree, error) {
	tree := &RedBlackTree{}
	if len(values) == 0 {
		return tree, nil
	}

	for _, v := range values {
		tree.Insert(v)
	}
	return tree, nil
}

// Height returns height of tree.
func (t *RedBlackTree) Height() int {
	return rbHeight(t.root)
}

// Size returns size of tree.
func (t *RedBlackTree) Size() int {
	return t.size
}

// Insert inserts a new node and keeps the tree an red-black tree. It should take O(logN) time.
func (t *RedBlackTree) Insert(value NodeValue) {
	ok, _, r := rbInsert(t.root, t.root, value, nil, true)
	if ok {
		t.size++
		t.root = r
		t.root.color = BLACK
	}
}

// Delete deletes specified node and keeps the tree an red-black tree. It should take O(logN) time.
func (t *RedBlackTree) Delete(value NodeValue) bool {
	r, b := rbDelete(t.root, t.root, value)
	if b {
		t.size--
		t.root = r
	}
	return b
}

// Validate returns whether the tree is a red-black tree which should be asserted true.
func (t *RedBlackTree) Validate() bool {
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

// Traverse traverse tree with specified order and call the function for each non-nil node.
func (t *RedBlackTree) Traverse(f func(value NodeValue), order TraverseOrder) {
	rbTraverse(t.root, f, order)
}

// String returns the inorder sequence and preorder sequence.
func (t *RedBlackTree) String() string {
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
func (t *RedBlackTree) PrettyPrint() string {
	if t.root == nil {
		return "empty"
	}
	type Ele struct {
		tab int
		n   *rbTreeNode
	}
	stack := make([]Ele, 0)
	stack = append(stack, Ele{tab: 0, n: t.root})
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

// Exist returns whether the value exists.
func (t *RedBlackTree) Exist(value NodeValue) bool {
	found, _, _ := rbSearch(value, t.root, nil)
	return found
}

func rbHeight(root *rbTreeNode) int {
	if root == nil {
		return 0
	}
	lh := rbHeight(root.left)
	rh := rbHeight(root.right)
	return int(math.Max(float64(lh), float64(rh))) + 1
}

// validateRb returns whether the sub-tree is valid red-black tree and the number of black nodes if true.
func validateRb(node *rbTreeNode) (bool, int) {
	if node == nil {
		return true, 1
	}
	if node.left != nil && node.left.value.Compare(node.value) >= 0 {
		return false, 0
	}
	if node.right != nil && node.right.value.Compare(node.value) <= 0 {
		return false, 0
	}
	if getNodeColor(node) == RED && (getNodeColor(node.left) == RED || getNodeColor(node.right) == RED) {
		return false, 0
	}
	lb, li := validateRb(node.left)
	rb, ri := validateRb(node.right)
	if lb && rb {
		if li != ri {
			return false, 0
		}
		if node.color == BLACK {
			return true, li + 1
		}
		return true, li
	}
	return false, 0
}

func rbInsert(root *rbTreeNode, node *rbTreeNode, value NodeValue, parent *rbTreeNode, left bool) (bool, bool, *rbTreeNode) {
	if root == nil {
		return true, true, &rbTreeNode{value: value, color: BLACK}
	}
	if node == nil {
		ele := &rbTreeNode{value: value, color: RED, parent: parent}
		if left {
			parent.left = ele
		} else {
			parent.right = ele
		}
		return true, false, root
	}
	if node.value.Compare(value) == 0 {
		return false, true, root
	}
	var ok, finish bool
	var r *rbTreeNode
	if node.value.Compare(value) > 0 {
		ok, finish, r = rbInsert(root, node.left, value, node, true)
	} else {
		ok, finish, r = rbInsert(root, node.right, value, node, false)
	}

	if !ok {
		return false, true, root
	} else {
		if finish || node.color == BLACK || parent == nil {
			return true, finish, r
		}
		uncle := parent.left
		if left {
			uncle = parent.right
		}
		if getNodeColor(uncle) == BLACK {
			parent.color = RED
			node.color = BLACK
			var rr *rbTreeNode
			if left {
				rr = rbRightRotate(root, parent)
			} else {
				rr = rbLeftRotate(root, parent)
			}
			return true, true, rr
		} else {
			node.color = BLACK
			parent.color = RED

			if uncle != nil {
				uncle.color = BLACK
			}
			return true, false, r
		}
	}
}

// returns finish, hasDeleted, new root
func rbDelete(root *rbTreeNode, node *rbTreeNode, value NodeValue) (*rbTreeNode, bool) {
	if node == nil {
		return root, false
	}
	if node.value.Compare(value) < 0 {
		return rbDelete(root, node.right, value)
	} else if node.value.Compare(value) > 0 {
		return rbDelete(root, node.left, value)
	}
	// case 0: node has no children and color is red: just delete
	if node.left == nil && node.right == nil {
		if node.color == RED {
			return detachChild(root, node), true
		} else {
			return rbDeleteNode(root, node), true
		}
	}
	// case 1: node has only one child
	if node.left != nil && node.right == nil || node.left == nil && node.right != nil {
		// if node has only one child, it must be red, just replace the node with its child and recolor to black
		if node.left != nil {
			node.left.color = BLACK
		} else {
			node.right.color = BLACK
		}
		return replaceWithChild(root, node), true
	}
	// case 2: node has two children
	s := rbFindSuccessor(node)
	node.value = s.value
	return rbDeleteNode(root, s), true
}

func rbTraverse(root *rbTreeNode, f func(value NodeValue), order TraverseOrder) {
	if root == nil {
		return
	}
	switch order {
	case PreOrder:
		f(root.value)
		rbTraverse(root.left, f, order)
		rbTraverse(root.right, f, order)
	case InOrder:
		rbTraverse(root.left, f, order)
		f(root.value)
		rbTraverse(root.right, f, order)
	case PostOrder:
		rbTraverse(root.left, f, order)
		rbTraverse(root.right, f, order)
		f(root.value)
	default:
		panic("unexpected order")
	}
}

func rbSearch(value NodeValue, root *rbTreeNode, lastVisited *rbTreeNode) (bool, *rbTreeNode, *rbTreeNode) {
	if root == nil {
		return false, nil, lastVisited
	}
	if root.value.Compare(value) == 0 {
		return true, root, lastVisited
	}
	if root.value.Compare(value) > 0 {
		return rbSearch(value, root.left, root)
	}
	return rbSearch(value, root.right, root)
}

func detachChild(root, node *rbTreeNode) *rbTreeNode {
	if node == nil {
		return root
	}
	if node.parent == nil {
		return nil
	}
	if node.parent.left == node {
		node.parent.left = nil
	} else {
		node.parent.right = nil
	}
	node.parent = nil
	return root
}

func replaceWithChild(root, node *rbTreeNode) *rbTreeNode {
	child := node.left
	if node.right != nil {
		child = node.right
	}
	if child != nil {
		child.parent = node.parent
	}
	if node.parent == nil {
		node.left = nil
		node.right = nil
		return child
	}
	if node.parent.left == node {
		node.parent.left = child
	} else {
		node.parent.right = child
	}
	node.parent = nil
	node.left = nil
	node.right = nil
	return root
}

func getBrother(node *rbTreeNode) *rbTreeNode {
	if node == nil {
		return nil
	}
	if node.parent == nil {
		return nil
	}
	if node.parent.left == node {
		return node.parent.right
	} else {
		return node.parent.left
	}
}

func getNephew(node *rbTreeNode, close bool) *rbTreeNode {
	brother := getBrother(node)
	if brother == nil {
		return nil
	}
	if node.parent.left == node {
		if close {
			return brother.left
		} else {
			return brother.right
		}
	} else {
		if close {
			return brother.right
		} else {
			return brother.left
		}
	}
}

func getNodeColor(node *rbTreeNode) nodeColor {
	if node == nil {
		return BLACK
	}
	return node.color
}

// left rotate tree right and return new root
func rbRightRotate(root *rbTreeNode, node *rbTreeNode) *rbTreeNode {
	l := node.left
	node.left = l.right
	if l.right != nil {
		l.right.parent = node
	}
	l.right = node
	l.parent = node.parent
	node.parent = l
	// case 0: parent is absent, meaning node is root
	if l.parent == nil {
		return l
	}
	// case 1: parent is present
	if l.parent.left == node {
		l.parent.left = l
	} else {
		l.parent.right = l
	}
	return root
}

// left rotate tree left and return new root
func rbLeftRotate(root *rbTreeNode, node *rbTreeNode) *rbTreeNode {
	r := node.right
	node.right = r.left
	if r.left != nil {
		r.left.parent = node
	}
	r.left = node
	r.parent = node.parent
	node.parent = r
	// case 0: parent is absent, meaning node is root
	if r.parent == nil {
		return r
	}
	// case 1: parent is present
	if r.parent.left == node {
		r.parent.left = r
	} else {
		r.parent.right = r
	}
	return root
}

func rbFindSuccessor(node *rbTreeNode) *rbTreeNode {
	q := node
	s := q.right
	for s.left != nil {
		q = s
		s = s.left
	}
	return s
}

func rbDeleteNode(root, node *rbTreeNode) *rbTreeNode {
	// case 0: ele has no children and it's color is red, just deleted
	if node.left == nil && node.right == nil && node.color == RED {
		return detachChild(root, node)
	}
	// case 1: ele has only right red child, just replace with child and recolor child to black
	if node.right != nil {
		node.right.color = BLACK
		return replaceWithChild(root, node)
	}
	// else: ele has no children and its color is black
	for q := node; q != root; {
		r1 := rbLeftRotate
		r2 := rbRightRotate
		if q.parent.left == q {
			r1 = rbRightRotate
			r2 = rbLeftRotate
		}
		// case 2.0: brother is black and close nephew is red
		if getNodeColor(getBrother(q)) == BLACK && getNodeColor(getNephew(q, true)) == RED {
			getBrother(q).color, getNephew(q, true).color = getNephew(q, true).color, getBrother(q).color
			root = r1(root, getBrother(q))
			continue
		}
		// case 2.1: brother is black and far nephew is red
		if getNodeColor(getBrother(q)) == BLACK && getNodeColor(getNephew(q, false)) == RED {
			q.parent.color, getBrother(q).color = getBrother(q).color, q.parent.color
			getNephew(q, false).color = BLACK
			r := r2(root, q.parent)
			return detachChild(r, node)
		}
		// case 2.2: brother is red
		if getNodeColor(getBrother(q)) == RED {
			q.parent.color = RED
			getBrother(q).color = BLACK
			root = r2(root, q.parent)
			continue
		}
		// case 2.3: parent is red and brother is black without children
		if getNodeColor(q.parent) == RED && getNodeColor(getBrother(q)) == BLACK {
			q.parent.color = BLACK
			getBrother(q).color = RED
			return detachChild(root, node)
		}
		// case 2.4: parent is black, brother has two black children (including nil)
		if getNodeColor(q.parent) == BLACK && getNodeColor(getBrother(q)) == BLACK {
			getBrother(q).color = RED
			// keep on
			q = q.parent
		}
	}
	return detachChild(root, node)
}

type nodeColor int

func (n nodeColor) String() string {
	if n == RED {
		return "red"
	} else {
		return "black"
	}
}

const (
	RED nodeColor = iota
	BLACK

	red   = "\033[31m"
	black = "\033[37m"
)
