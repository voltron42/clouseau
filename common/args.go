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

func (a *Args) PopLast() *Arg {
	size := len(*a)
	if size == 0 {
		return nil
	}
	last := a.Get(size - 1)
	*a = (*a)[:size-1]
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

func (a *Args) Subset(start, end int) *Args {
	out := *a
	if start != 0 || end != -1 {
		if start == 0 {
			out = out[:end]
		} else if end == -1 {
			out = out[start:]
		} else {
			out = out[start:end]
		}
	}
	return &out
}

func (a *Args) Some(fn func(item *Arg, index int) bool) *Args {
	out := Args{}
	size := a.Len()
	for x := 0; x < size; x++ {
		value := a.Get(x)
		if fn(value, x) {
			out = append(out, value.Elem())
		}
	}
	return &out
}

func (a *Args) Map(fn func(item *Arg, index int) *Arg) *Args {
	out := Args{}
	size := a.Len()
	for x := 0; x < size; x++ {
		out = append(out, fn(a.Get(x), x).Elem())
	}
	return &out
}

func (a *Args) Each(fn func(item *Arg, index int)) {
	size := a.Len()
	for x := 0; x < size; x++ {
		fn(a.Get(x), x)
	}
}

func (a *Args) Len() int {
	return len(*a)
}

func (a *Args) Get(index int) *Arg {
	if len(*a) <= index {
		return &Arg{}
	}
	return &Arg{(*a)[index]}
}

type Arg struct {
	inner interface{}
}

func NewArg(elem interface{}) *Arg {
	return &Arg{elem}
}

func (a *Arg) Elem() interface{} {
	return a.inner
}

func (a *Arg) TypeOf() reflect.Type {
	return reflect.TypeOf(a.inner)
}

func (a *Arg) ValueOf() reflect.Value {
	return reflect.ValueOf(a.inner)
}

func (a *Arg) Bool() bool {
	return a.ValueOf().Bool()
}

func (a *Arg) String() string {
	return fmt.Sprintf("%v", a.inner)
}

func (a *Arg) Strings() []string {
	return castAsStrings(a.inner)
}

func (a *Arg) Error() error {
	return castAsError(a.inner)
}

func (a *Arg) Bytes() []byte {
	return castAsBytes(a.inner)
}

func (a *Arg) Byte() byte {
	return tribble.NewTribble(a.inner).Byte()
}

func (a *Arg) Float32() float32 {
	return tribble.NewTribble(a.inner).Float32()
}

func (a *Arg) Float64() float64 {
	return tribble.NewTribble(a.inner).Float64()
}

func (a *Arg) Int() int {
	return tribble.NewTribble(a.inner).Int()
}

func (a *Arg) Int8() int8 {
	return tribble.NewTribble(a.inner).Int8()
}

func (a *Arg) Int16() int16 {
	return tribble.NewTribble(a.inner).Int16()
}

func (a *Arg) Int32() int32 {
	return tribble.NewTribble(a.inner).Int32()
}

func (a *Arg) Int64() int64 {
	return tribble.NewTribble(a.inner).Int64()
}

func (a *Arg) UInt() uint {
	return tribble.NewTribble(a.inner).UInt()
}

func (a *Arg) UInt16() uint16 {
	return tribble.NewTribble(a.inner).UInt16()
}

func (a *Arg) UInt32() uint32 {
	return tribble.NewTribble(a.inner).UInt32()
}

func (a *Arg) UInt64() uint64 {
	return tribble.NewTribble(a.inner).UInt64()
}

func (a *Arg) UIntPtr() uintptr {
	return tribble.NewTribble(a.inner).UIntPtr()
}
