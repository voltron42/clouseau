package mockband

import (
	"../common"
	"reflect"
)

type Mock struct {
	calls map[string]*callList
}

func NewMock() *Mock {
	return &Mock{map[string]*callList{}}
}

func (m *Mock) Called(name string, params ...interface{}) *common.Args {
	call := m.getCall(name, params)
	if call == nil {
		panic("Function with param signature not found: " + name)
	} else {
		args := common.Args(params)
		return call.exec(&args)
	}
}

func (m *Mock) CalledVarArg(name string, params ...interface{}) *common.Args {
	if len(params) == 0 {
		return m.Called(name)
	} else {
		args := common.Args(params)
		last := reflect.ValueOf(args.PopLast())
		err := args.AddAll(last)
		if err != nil {
			args.Add(last)
		}
		return m.Called(name, args...)
	}
}

func (m *Mock) When(name string, params ...interface{}) *call {
	list, ok := m.calls[name]
	if !ok {
		list = &callList{}
		m.calls[name] = list
	}
	return list.createCall(params)
}

func (m *Mock) GetCalls(name string, params ...interface{}) *results {
	call := m.getCall(name, params)
	if call == nil {
		panic("Function not found: " + name)
	}
	return &call.results
}

func (m *Mock) HasCalled(name string, params ...interface{}) *metric {
	call := m.getCall(name, params)
	if call == nil {
		panic("Function not found: " + name)
	}
	return &metric{call.results.count()}
}

func (m *Mock) getCall(name string, params []interface{}) *call {
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

func (r *results) GetResults(index int) *common.Args {
	return r.list[index].results
}

func (r *results) GetParams(index int) *common.Args {
	return r.list[index].params
}

func (r *results) add(result result) {
	r.list = append(r.list, result)
}

func (r *results) count() int {
	return len(r.list)
}

type result struct {
	params  *common.Args
	results *common.Args
}

type call struct {
	list    []func(args *common.Args) *common.Args
	index   int
	results results
}

func newCall() *call {
	return &call{
		[]func(args *common.Args) *common.Args{},
		0,
		results{},
	}
}

func (c *call) exec(args *common.Args) *common.Args {
	out := c.list[c.index](args)
	c.results.add(result{
		params:  args,
		results: out,
	})
	c.index = (c.index + 1) % len(c.list)
	return out
}

func (c *call) Return(params ...interface{}) *call {
	return c.Then(func(args *common.Args) *common.Args {
		out := common.Args(params)
		return &out
	})
}

func (c *call) Inject(value interface{}, index int, params ...interface{}) *call {
	return c.Then(func(args *common.Args) *common.Args {
		pointer := reflect.ValueOf((*args)[index])
		pointer.Elem().Set(reflect.ValueOf(value))
		out := common.Args(params)
		return &out
	})
}

func (c *call) Panic(err interface{}) *call {
	return c.Then(func(args *common.Args) *common.Args {
		panic(err)
	})
}

func (c *call) Then(fn func(args *common.Args) *common.Args) *call {
	c.list = append(c.list, fn)
	return c
}

type callListItem struct {
	params common.Args
	call   *call
}

type callList struct {
	list []callListItem
}

func (c *callList) getCall(params []interface{}) *call {
	for _, item := range c.list {
		args := common.Args(params)
		if item.params.Matches(&args) {
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
	me = newCall()
	c.list = append(c.list, callListItem{
		params: common.Args(params),
		call:   me,
	})
	return me
}
