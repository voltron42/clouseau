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

func (m *Mock) GetCall(name string, params ...interface{}) *results {
	call := m.getCall(name, params)
	if call == nil {
		reckon.Fail(errors.New("Function not found: " + name))
	}
	return &metric{call.results.count()}
}

func (m *Mock) HasCalled(name string, params ...interface{}) *metric {
	call := m.getCall(name, params)
	if call == nil {
		reckon.Fail(errors.New("Function not found: " + name))
	}
	return &call.results
}

func (m *Mock) getCall(name string, params ...interface{}) *call {
	list, ok := m.calls[name]
	if !ok {
		return nil
	}
	return list.getCall(params)
}

type metric struct {
	count int
}

func (m *metric) Times(times int) bool {
	return m.count == times
}

func (m *metric) Once() bool {
	return m.Times(1)
}

func (m *metric) Twice() bool {
	return m.Times(2)
}

type results struct {
	list []result
}

func (r *results) GetResults(index int) *Args {
	return r.list[index].results
}

func (r *results) GetParams(index int) *Args {
	return r.list[index].params
}

func (r *results) add(result result) {
	r.list = append(r.list, result)
}

func (r *results) count() int {
	return len(r.list)
}

type result struct {
	params  *Args
	results *Args
}

type call struct {
	list    []func(args *Args) *Args
	index   int
	results results
}

func (c *call) exec(args *Args) *Args {
	out := c.list[c.index](args)
	c.results.add(result{
		params:  args,
		results: out,
	})
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

func Any() interface{} {
	return any
}

var any interface{} = struct{}{}
