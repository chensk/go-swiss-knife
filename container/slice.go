package container

import (
	"github.com/chensk/go-swiss-knife/assert"
	"reflect"
)

// RemoveFromList removes element from list and return the result list with original list not modified.
// Note it removes only the first element found.
func RemoveFromList(list interface{}, element interface{}) (interface{}, error) {
	// assert element type matches list
	t := reflect.TypeOf(list)
	if err := assert.Assert(t.Kind() == reflect.Slice, "first argument is not a slice"); err != nil {
		return nil, err
	}
	if err := assert.Assert(t.Elem() == reflect.TypeOf(element), "slice element type doesn't match target element type"); err != nil {
		return nil, err
	}
	v := reflect.ValueOf(list)
	if v.Len() == 0 {
		return list, nil
	}
	var i int
	for i = 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == element {
			break
		}
	}
	// not found
	if i >= v.Len() {
		return list, nil
	}
	res := reflect.MakeSlice(t, v.Len()-1, v.Len())
	for ii := 0; ii < v.Len(); ii++ {
		if ii < i {
			res.Index(ii).Set(v.Index(ii))
		} else if ii == i {
			continue
		} else {
			res.Index(ii - 1).Set(v.Index(ii))
		}
	}
	return res.Interface(), nil
}
