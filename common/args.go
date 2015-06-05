package common

import (
	"../../tribble"
	"errors"
	"fmt"
	"reflect"
)

func Any() interface{} {
	return any
}

var any interface{} = struct{}{}

type Args []interface{}

func (a *Args) PopLast() interface{} {
	size := len(*a)
	last := (*a)[size-1]
	*a = (*a)[:size-2]
	return last
}

func (a *Args) AddAll(value reflect.Value) error {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return errors.New("not slice or array")
	}
	count := value.Len()
	for x := 0; x < count; x++ {
		*a = append(*a, value.Index(x).Interface())
	}
	return nil
}

func (a *Args) Add(value reflect.Value) {
	*a = append(*a, value.Interface())
}

func (a *Args) Matches(args *Args) bool {
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

func (a *Args) Error(index int) error {
	err, ok := a.Get(index).(error)
	if !ok {
		panic(fmt.Sprintf("Arg %v cannot be cast to error.", index))
	}
	return err
}

func (a *Args) Byte(index int) byte {
	return tribble.NewTribble(a.Get(index)).Byte()
}

func (a *Args) Float32(index int) float32 {
	return tribble.NewTribble(a.Get(index)).Float32()
}

func (a *Args) Float64(index int) float64 {
	return tribble.NewTribble(a.Get(index)).Float64()
}

func (a *Args) Int(index int) int {
	return tribble.NewTribble(a.Get(index)).Int()
}

func (a *Args) Int8(index int) int8 {
	return tribble.NewTribble(a.Get(index)).Int8()
}

func (a *Args) Int16(index int) int16 {
	return tribble.NewTribble(a.Get(index)).Int16()
}

func (a *Args) Int32(index int) int32 {
	return tribble.NewTribble(a.Get(index)).Int32()
}

func (a *Args) Int64(index int) int64 {
	return tribble.NewTribble(a.Get(index)).Int64()
}

func (a *Args) UInt(index int) uint {
	return tribble.NewTribble(a.Get(index)).UInt()
}

func (a *Args) UInt16(index int) uint16 {
	return tribble.NewTribble(a.Get(index)).UInt16()
}

func (a *Args) UInt32(index int) uint32 {
	return tribble.NewTribble(a.Get(index)).UInt32()
}

func (a *Args) UInt64(index int) uint64 {
	return tribble.NewTribble(a.Get(index)).UInt64()
}

func (a *Args) UIntPtr(index int) uintptr {
	return tribble.NewTribble(a.Get(index)).UIntPtr()
}
