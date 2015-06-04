package reckon

import (
	"reflect"
)

func That(actual interface{}) *reckoning {
	return &reckoning{
		newIs(actual),
		newDoes(actual),
		newHas(actual),
	}
}

type reckoning struct {
	Is   *is
	Does *does
	Has  *has
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
		objCompare: &objCompare{actual: actual, state: true},
		Not:        &objCompare{actual: actual, state: false},
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

func (o *owner) Property(name string, values ...interface{}) {
	o.getProp("has property", name, values)
}

func (o *owner) Key(name string, values ...interface{}) {
	o.getProp("has key", name, values)
}

func (o *owner) DeepProperty(name string, values ...interface{}) {
	o.getProp("has deep property", name, values)
}

func (o *owner) getProp(fn, name string, values []interface{}) {
	params := append([]interface{}{}, o.actual, name)
	for _, value := range values {
		params = append(params, value)
	}
	expectations.check("has deep property", o.state, params)
}

type length struct {
	*numberCompare
	Not *numberCompare
}

func newLength(actual interface{}) *length {
	return &length{
		numberCompare: &numberCompare{actual, true},
		Not:           &numberCompare{actual, false},
	}
}

type listing struct {
	actual interface{}
	state  bool
}

func (l *listing) Keys(names ...string) {
	l.getList("has keys", names)
}

func (l *listing) Properties(names ...string) {
	l.getList("has properties", names)
}

func (l *listing) getList(fn string, names []string) {
	params := append([]interface{}{}, l.actual, l.state)
	for _, name := range names {
		params = append(params, name)
	}
	expectations.check(fn, true, params)
}
