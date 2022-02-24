package container

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	t1 := time.Now()
	//tree := createRandomAvlTree(1000, true)
	tree := createRandomRbTree(10000)
	t.Logf("create: %v\n", time.Now().Sub(t1))
	t.Logf("height: %d, size: %d, validate: %t\n", tree.Height(), tree.Size(), tree.Validate())
	t1 = time.Now()
	t.Logf("stats: %v\n", time.Now().Sub(t1))
	if !tree.Validate() {
		t.Fatalf("tree invalid: %v\n", tree.PrettyPrint())
	}
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
	}
	t.Logf("tree: %s\n", tree.PrettyPrint())
}

func TestPop(t *testing.T) {
	tree := createRandomRbTree(0)
	t.Logf("height: %d, size: %d, validate: %t\n", tree.Height(), tree.Size(), tree.Validate())
	for i := 0; i < 100; i++ {
		v := tree.PopMin(Element(1000), true)
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
	rand.Seed(time.Now().UnixNano())
	input := make([]NodeValue, size)
	ii := make([]int, 0, size)
	for i := 0; i < size; i++ {
		ii = append(ii, i)
	}
	rand.Shuffle(len(ii), func(i, j int) {
		ii[i], ii[j] = ii[j], ii[i]
	})
	for i := 0; i < size; i++ {
		input[i] = Element(ii[i])
	}
	fmt.Printf("inputs: %v\n", input)
	tree, err := NewBinarySearchTree(input)
	if err != nil {
		panic(err)
	}
	return tree
}

func createRandomRbTree(size int) *RedBlackTree {
	rand.Seed(time.Now().UnixNano())
	input := make([]NodeValue, size)
	ii := make([]int, 0, size)
	for i := 0; i < size; i++ {
		ii = append(ii, i)
	}
	rand.Shuffle(len(ii), func(i, j int) {
		ii[i], ii[j] = ii[j], ii[i]
	})
	for i := 0; i < size; i++ {
		input[i] = Element(ii[i])
	}
	tree := NewRedBlackTree(input)
	return tree
}

func TestIterator(t *testing.T) {
	tree := createRandomRbTree(1000)
	iter, err := tree.Iterator(nil, nil, InOrder)
	if err != nil {
		t.Fatal(err)
	}
	for ele := iter.Next(); ele != nil; ele = iter.Next() {
		t.Logf("get %v\n", ele)
	}
}

type Element int

func (e Element) Compare(value NodeValue) int {
	return int(e - value.(Element))
}

func (e Element) String() string {
	return strconv.Itoa(int(e))
}

// coral obj, 40B in size
type CoralObj struct {
	// for lru: last accessed time stamp
	// todo: for lfu
	collectFlag int32
	otype       Otype
	// if otype is OBJ_STRING_RAW, ptr is byte slice storing the value itself
	// if otype is OBJ_STRING_PERSISTENCE, ptr is pointer to CoralStringPersistence
	ptr interface{}
	key string
}

type Otype int8

const (
	OBJ_STRING_RAW Otype = 1 << iota
	OBJ_STRING_PERSISTENCE
)

func (o *CoralObj) Compare(value NodeValue) int {
	c1 := value.(*CoralObj)
	ret := int(o.collectFlag - c1.collectFlag)
	if ret != 0 {
		return ret
	}
	return strings.Compare(o.key, c1.key)
}

func (o *CoralObj) String() string {
	return o.key
}

func (o *CoralObj) Size() int64 {
	switch o.otype {
	case OBJ_STRING_PERSISTENCE:
		return o.ptr.(*CoralStringPersistence).size
	case OBJ_STRING_RAW:
		return int64(len(o.ptr.([]byte)))
	default:
		panic("unexpected object type")
	}
}

type CoralStringPersistence struct {
	path string
	size int64
}
