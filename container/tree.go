package container

import "fmt"

type BinarySearchTree interface {
	// Height returns height of tree.
	Height() int
	// Size returns size of tree.
	Size() int
	Insert(value NodeValue)
	Delete(value NodeValue) bool
	Validate() bool
	Traverse(f func(value NodeValue) bool, order TraverseOrder)
	String() string
	PrettyPrint() string
	Exist(value NodeValue) bool
}

// NodeValue represents the element type storing in tree.
type NodeValue interface {
	// Compare method compares with another NodeValue and return 1 if greater, 0 if equal, -1 if less.
	Compare(value NodeValue) int

	// Stringer whose String() method would return formatted string.
	fmt.Stringer
}

type TreeIterator struct {
	ch     chan NodeValue
	closed chan struct{}
}

func (iter TreeIterator) Next() NodeValue {
	select {
	case v, ok := <-iter.ch:
		if !ok {
			return nil
		}
		return v
	case <-iter.closed:
		return nil
	}
}

func (iter TreeIterator) Close() {
	select {
	case <-iter.closed:
		return
	default:
		close(iter.closed)
		close(iter.ch)
	}
}

func (iter TreeIterator) HasNext() bool {
	select {
	case <-iter.closed:
		return false
	default:
		return true
	}
}

type TraverseOrder int

const (
	// PreOrder parent-left-right
	PreOrder TraverseOrder = iota
	// InOrder left-parent-right
	InOrder
	// PostOrder left-right-parent
	PostOrder
	// ReversedOrder right-parent-left
	ReversedOrder
)
