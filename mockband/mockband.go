package mockband

import (
	"../reckon"
	"errors"
	"reflect"
)

type Mock struct {
	calls map[string]callList
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Called(name string, params ...interface{}) *Args {
	call := m.getCall(name, params)
	if call == nil {
		reckon.Fail(errors.New("Function with param signature not found: " + name))
		return nil
	} else {
		args := Args(params)
		return call.exec(&args)
	}
}

func (m *Mock) CalledVarArg(name string, params ...interface{}) *Args {
	if len(params) == 0 {
		return m.Called(name)
	} else {
		args := Args(params)
		last := reflect.ValueOf(args.popLast())
		err := args.addAll(last)
		if err != nil {
			args.add(last)
		}
		return m.Called(name, args...)
	}
}

func (m *Mock) When(name string, params ...interface{}) *call {
	list, ok := m.calls[name]
	if !ok {
		list = callList{}
		m.calls[name] = list
	}
	return list.createCall(params)
}

func (m *Mock) GetCall(name string, params ...interface{}) *call {
	call := m.getCall(name, params)
	if call == nil {
		reckon.Fail(errors.New("Function not found: " + name))
	}
	return call
}

func (m *Mock) getCall(name string, params ...interface{}) *call {
	list, ok := m.calls[name]
	if !ok {
		return nil
	}
	return list.getCall(params)
}

type call struct {
	list  []func(args *Args) *Args
	index int
}

func (c *call) exec(args *Args) *Args {
	out := c.list[c.index](args)
	c.index = (c.index + 1) % len(c.list)
	return out
}

func (c *call) Return(params ...interface{}) *call {
	return c.Then(func(args *Args) *Args {
		out := Args(params)
		return &out
	})
}

func (c *call) Inject(value interface{}, index int, params ...interface{}) *call {
	return c.Then(func(args *Args) *Args {
		pointer := reflect.ValueOf((*args)[index])
		pointer.Elem().Set(reflect.ValueOf(value))
		out := Args(params)
		return &out
	})
}

func (c *call) Panic(err interface{}) *call {
	return c.Then(func(args *Args) *Args {
		panic(err)
	})
}

func (c *call) Then(fn func(args *Args) *Args) *call {
	c.list = append(c.list, fn)
	return c
}

type callList struct {
	list []struct {
		params Args
		call   *call
	}
}

func (c *callList) getCall(params []interface{}) *call {
	for _, item := range c.list {
		args := Args(params)
		if item.params.matches(&args) {
			return item.call
		}
	}
	return nil
}

func (c *callList) createCall(params []interface{}) *call {
	me := c.getCall(params)
	if me != nil {
		return me
	}
	me = &call{}
	c.list = append(c.list, struct {
		params Args
		call   *call
	}{
		params: Args(params),
		call:   me,
	})
	return me
}

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
	// TODO --
	return false
}

var Any interface{} = struct{}{}
