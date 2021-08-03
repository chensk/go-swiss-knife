package container

import "testing"

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
	}

	for _, arg := range args {
		r, e := RemoveFromList(arg.sli, arg.ele)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(r)
	}
}
