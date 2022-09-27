package container

import (
	"errors"
	"fmt"
)

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

func CreateTreeIterator(tree BinarySearchTree, from NodeValue, to NodeValue, order TraverseOrder) (*TreeIterator, error) {
	if from != nil && to != nil && from.Compare(to) > 0 {
		return nil, errors.New("invalid bound")
	}
	ch := make(chan NodeValue)
	closed := make(chan struct{})
	go func() {
		tree.Traverse(func(value NodeValue) bool {
			if (from == nil || value.Compare(from) >= 0) && (to == nil || value.Compare(to) <= 0) {
				// double check closed channel to prevent nil panic
				select {
				case <-closed:
					return false
				default:
				}

				select {
				case <-closed:
					return false
				case ch <- value:
				}
			}
			switch order {
			case PreOrder, PostOrder:
				return true
			case InOrder:
				return to == nil || value.Compare(to) <= 0
			case ReversedOrder:
				return from == nil || value.Compare(from) >= 0
			default:
				return false
			}
		}, order)
		select {
		case <-closed:
			return
		default:
			close(closed)
			close(ch)
		}
	}()
	return &TreeIterator{
		ch:     ch,
		closed: closed,
	}, nil
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
