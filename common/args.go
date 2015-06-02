package common

import (
	"errors"
	"fmt"
	"reflect"
)

func Any() interface{} {
	return any
}

var any interface{} = struct{}{}

type Args []interface{}

func (a *Args) popLast() interface{} {
	size := len(*a)
	last := (*a)[size-1]
	*a = (*a)[:size-2]
	return last
}

func (a *Args) addAll(value reflect.Value) error {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return errors.New("not slice or array")
	}
	count := value.Len()
	for x := 0; x < count; x++ {
		*a = append(*a, value.Index(x).Interface())
	}
	return nil
}

func (a *Args) add(value reflect.Value) {
	*a = append(*a, value.Interface())
}

func (a *Args) matches(args *Args) bool {
	if len(*args) > len(*a) {
		return false
	}
	count := len(*args)
	for x := 0; x < count; x++ {
		if !reflect.DeepEqual((*a)[x], (*args)[x]) && !reflect.DeepEqual((*a)[x], any) {
			return false
		}
	}
	return true
}

func (a *Args) Len() int {
	return len(*a)
}

func (a *Args) Get(index int) interface{} {
	if len(*a) <= index {
		return nil
	}
	return (*a)[index]
}

func (a *Args) TypeOf(index int) reflect.Type {
	return reflect.TypeOf(a.Get(index))
}

func (a *Args) ValueOf(index int) reflect.Value {
	return reflect.ValueOf(a.Get(index))
}

func (a *Args) Bool(index int) bool {
	return a.ValueOf(index).Bool()
}

func (a *Args) String(index int) string {
	return fmt.Sprintf("%v", a.Get(index))
}
