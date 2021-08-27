package container

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	args := []struct {
		sli interface{}
		ele interface{}
	}{
		{
			sli: []int{1, 2, 3, 4, 5, 6},
			ele: 0,
		},
		{
			sli: []int{1, 2, 3, 4, 5, 6},
			ele: 3,
		},
		{
			sli: []int{1, 2, 3, 4, 5, 6},
			ele: 6,
		},
		{
			sli: []string{"abc", "def", "ghi"},
			ele: "def",
		},
		{
			sli: []string{"abc", "def", "ghi"},
			ele: "deff",
		},
		{
			sli: nil,
			ele: "test",
		},
		{
			sli: "abcd",
			ele: "a",
		},
	}

	for _, arg := range args {
		s := Of(arg.sli)
		s = s.Remove(arg.ele)
		slice, err := s.Element()
		if err != nil {
			t.Error(err)
		} else {
			t.Log(slice)
		}
	}
}

func TestBinarySearchTree(t *testing.T) {
	tree := createRandomTree(10000)
	t.Logf("height: %d, size: %d, validate: %t\n", tree.Height(), tree.Size(), tree.ValidateAvl())

	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)
	for tree.Size() > 0 {
		inOrders := make([]NodeValue, 0, tree.Size())
		tree.Traverse(func(value NodeValue) {
			inOrders = append(inOrders, value)
		}, InOrder)
		if len(inOrders) != tree.Size() {
			panic("unexpected")
		}
		i := inOrders[rnd.Intn(tree.Size())]
		tree.Delete(i)
		if !tree.ValidateAvl() {
			t.Fatal("tree is not valid")
		}
	}
}

func createRandomTree(size int) *BalancedBinarySearchTree {
	input := make([]NodeValue, size)
	for i := 0; i < size; i++ {
		input[i] = Element(i)
	}
	tree, err := NewBinarySearchTree(input)
	if err != nil {
		panic(err)
	}
	return tree
}

type Element int

func (e Element) Compare(value NodeValue) int {
	return int(e - value.(Element))
}

func (e Element) String() string {
	return strconv.Itoa(int(e))
}
