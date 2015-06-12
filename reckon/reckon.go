package reckon

import (
	"reflect"
)

func That(actual interface{}) *reckoning {
	return &reckoning{
		newIs(actual),
		newDoes(actual),
		newHas(actual),
		newWill(actual),
	}
}

type reckoning struct {
	Is   *is
	Does *does
	Has  *has
	Will *will
}

type does struct {
	*objContains
	Not *objContains
}

func newDoes(actual interface{}) *does {
	return &does{
		objContains: &objContains{actual: actual, state: true},
		Not:         &objContains{actual: actual, state: false},
	}
}

type objContains struct {
	actual interface{}
	state  bool
}

func (o *objContains) Exist() {
	expectations.check("exists", o.state, o.actual)
}

func (o *objContains) Match(regex string) {
	expectations.check("matches", o.state, o.actual, regex)
}

func (o *objContains) Contain(needle string) {
	expectations.check("contains", o.state, o.actual, needle)
}

// Contains

type is struct {
	*objCompare
	Not *objCompare
}

func newIs(actual interface{}) *is {
	return &is{
		objCompare: &objCompare{numberCompare: &numberCompare{actual, true}, actual: actual, state: true},
		Not:        &objCompare{numberCompare: &numberCompare{actual, false}, actual: actual, state: false},
	}
}

type objCompare struct {
	*numberCompare
	actual interface{}
	state  bool
}

func (o *objCompare) EqualTo(expected interface{}) {
	expectations.check("equals", o.state, o.actual, expected)
}

func (o *objCompare) Nil() {
	o.EqualTo(nil)
}

func (o *objCompare) True() {
	o.EqualTo(true)
}

func (o *objCompare) False() {
	o.EqualTo(false)
}

func (o *objCompare) Zero() {
	expectations.check("is zero", o.state, o.actual)
}

func (o *objCompare) A(kind reflect.Kind) {
	expectations.check("is a", o.state, o.actual, kind)
}

func (o *objCompare) AnInstanceOf(t reflect.Type) {
	expectations.check("instance of", o.state, o.actual, t)
}

type numberCompare struct {
	actual interface{}
	state  bool
}

func (n *numberCompare) GreaterThan(bound float64) {
	expectations.check("greater than", n.state, n.actual, bound)
}

func (n *numberCompare) LessThan(bound float64) {
	expectations.check("less than", n.state, n.actual, bound)
}

func (n *numberCompare) Within(low float64, high float64) {
	expectations.check("within", n.state, n.actual, low, high)
}

type has struct {
	*owner
	No     *owner
	Any    *listing
	All    *listing
	Length *length
}

func newHas(actual interface{}) *has {
	return &has{
		owner:  &owner{actual, true},
		No:     &owner{actual, false},
		Any:    &listing{actual, true},
		All:    &listing{actual, false},
		Length: newLength(actual),
	}
}

type owner struct {
	actual interface{}
	state  bool
}

func (o *owner) Property(name interface{}, values ...interface{}) {
	o.getProp("has property", name, values)
}

func (o *owner) DeepProperty(name string, values ...interface{}) {
	o.getProp("has deep property", name, values)
}

func (o *owner) getProp(fn string, name interface{}, values []interface{}) {
	params := append([]interface{}{}, o.actual, name)
	for _, value := range values {
		params = append(params, value)
	}
	expectations.check(fn, o.state, params...)
}

type length struct {
	*numberCompare
	Not *numberCompare
}

func newLength(actual interface{}) *length {
	temp := actual
	value := reflect.ValueOf(actual)
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		temp = value.Len()
	}
	return &length{
		numberCompare: &numberCompare{temp, true},
		Not:           &numberCompare{temp, false},
	}
}

type listing struct {
	actual interface{}
	state  bool
}

func (l *listing) Keys(names ...interface{}) {
	l.getList("has keys", names)
}

func (l *listing) Properties(names ...interface{}) {
	l.getList("has properties", names)
}

func (l *listing) getList(fn string, names []interface{}) {
	params := append([]interface{}{}, l.actual, l.state)
	params = append(params, names...)
	expectations.check(fn, true, params...)
}

type will struct {
	*panicCheck
	Not *panicCheck
}

func newWill(actual interface{}) *will {
	return &will{
		panicCheck: &panicCheck{actual, true},
		Not:        &panicCheck{actual, false},
	}
}

type panicCheck struct {
	actual interface{}
	state  bool
}

func (p *panicCheck) Panic() {
	expectations.check("panics", p.state, p.actual)
}

func (p *panicCheck) PanicWith(message interface{}) {
	expectations.check("panics with message", p.state, p.actual, message)
}
