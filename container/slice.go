package container

import (
	"github.com/chensk/go-swiss-knife/assert"
	"reflect"
)

type Slice struct {
	_slice       interface{}
	_elementType reflect.Type
	_err         error
}

// Of constructs a Slice object from a go slice, if argument is not a valid slice, a specific Slice Object would be returned
// which any operations on it would failed.
func Of(slice interface{}) *Slice {
	if err := assert.Assert(slice != nil, "slice is nil"); err != nil {
		return &Slice{_err: err}
	}
	t := reflect.TypeOf(slice)
	if err := assert.Assert(t.Kind() == reflect.Slice, "argument is not a slice"); err != nil {
		return &Slice{_err: err}
	}
	return &Slice{
		_slice:       slice,
		_elementType: t.Elem(),
		_err:         nil,
	}
}

func (s *Slice) Error() error {
	return s._err
}

// Remove removes element from list and return the result list with original list not modified.
// Note it removes only the first element found.
func (s *Slice) Remove(element interface{}) *Slice {
	if s._err != nil {
		return s
	}
	if err := assert.Assert(s._elementType == reflect.TypeOf(element), "slice element type doesn't match target element type"); err != nil {
		s._err = err
		return s
	}
	v := reflect.ValueOf(s._slice)
	if v.Len() == 0 {
		return s
	}
	var i int
	for i = 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == element {
			break
		}
	}
	// not found
	if i >= v.Len() {
		return s
	}
	res := reflect.MakeSlice(reflect.TypeOf(s._slice), v.Len()-1, v.Len())
	for ii := 0; ii < v.Len(); ii++ {
		if ii < i {
			res.Index(ii).Set(v.Index(ii))
		} else if ii == i {
			continue
		} else {
			res.Index(ii - 1).Set(v.Index(ii))
		}
	}
	s._slice = res.Interface()
	return s
}

// ListExist returns whether element exists in list.
// Note it returns false if any type errors happens.
func (s *Slice) ListExist(element interface{}) bool {
	if s._err != nil {
		return false
	}
	if s._elementType != reflect.TypeOf(element) {
		return false
	}
	v := reflect.ValueOf(s._slice)
	if v.Len() == 0 {
		return false
	}
	var i int
	for i = 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == element {
			return true
		}
	}
	return false
}

// Element gets the underlying slice.
func (s *Slice) Element() (interface{}, error) {
	if s._err != nil {
		return nil, s._err
	}
	return s._slice, nil
}
