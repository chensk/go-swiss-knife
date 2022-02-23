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

func TestTree(t *testing.T) {
	//tree := createRandomAvlTree(1000, true)
	tree := createRandomRbTree(10000, true)
	t.Logf("height: %d, size: %d, validate: %t\n", tree.Height(), tree.Size(), tree.Validate())
	//t.Logf("tree: %s\n", tree.PrettyPrint())
	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)
	for tree.Size() > 0 {
		inOrders := make([]NodeValue, 0, tree.Size())
		tree.Traverse(func(value NodeValue) bool {
			inOrders = append(inOrders, value)
			return true
		}, InOrder)
		if len(inOrders) != tree.Size() {
			panic("unexpected")
		}
		i := inOrders[rnd.Intn(tree.Size())]
		tree.Delete(i)
		if !tree.Validate() {
			t.Fatalf("after deleted %d, tree is not valid, tree: %s\n", i, tree.PrettyPrint())
		}
		t.Logf("succesfully deleted %d, tree: %d\n", i, tree.Size())
	}
	t.Logf("tree: %s\n", tree.PrettyPrint())
}

func TestPop(t *testing.T) {
	tree := createRandomRbTree(1000, true)
	t.Logf("height: %d, size: %d, validate: %t\n", tree.Height(), tree.Size(), tree.Validate())
	for i := 0; i < 100; i++ {
		v := tree.Pop(func(value NodeValue) bool {
			return value.(Element).Compare(Element(50)) > 0 && value.(Element).Compare(Element(450)) < 0
		}, true)
		t.Logf("pop %v\n", v)
		if !tree.Validate() {
			t.Fatalf("after deleted %d, tree is not valid, tree: %s\n", i, tree.PrettyPrint())
		}
	}
	t.Logf("tree: %s\n size: %d\n", tree.PrettyPrint(), tree.Size())
}

func TestExpiringSet(t *testing.T) {
	s := NewExpiringSet(1 * time.Second)
	s.Add("123")
	time.Sleep(200 * time.Millisecond)
	s.Add("456")
	time.Sleep(200 * time.Millisecond)
	s.Add("789")

	t.Logf("exists: %t %t %t", s.Exists("123"), s.Exists("456"), s.Exists("789"))
	time.Sleep(700 * time.Millisecond)
	t.Logf("exists: %t %t %t", s.Exists("123"), s.Exists("456"), s.Exists("789"))
	time.Sleep(200 * time.Millisecond)
	t.Logf("exists: %t %t %t", s.Exists("123"), s.Exists("456"), s.Exists("789"))
	time.Sleep(200 * time.Millisecond)
	t.Logf("exists: %t %t %t", s.Exists("123"), s.Exists("456"), s.Exists("789"))
}

func createRandomAvlTree(size int, asc bool) *BalancedBinarySearchTree {
	input := make([]NodeValue, size)
	for i := 0; i < size; i++ {
		if asc {
			input[i] = Element(i)
		} else {
			input[i] = Element(size - i)
		}
	}
	tree, err := NewBinarySearchTree(input)
	if err != nil {
		panic(err)
	}
	return tree
}

func createRandomRbTree(size int, asc bool) *RedBlackTree {
	input := make([]NodeValue, size)
	for i := 0; i < size; i++ {
		if asc {
			input[i] = Element(i)
		} else {
			input[i] = Element(size - i)
		}
	}
	tree, err := NewRedBlackTree(input)
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
