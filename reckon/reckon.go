package reckon

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

func (o *objContains) Match(expected interface{}) {
	expectations.check("matches", o.state, o.actual, expected)
}

func (o *objContains) Contain(expected interface{}) {
	expectations.check("contains", o.state, o.actual, expected)
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
}

/*
Zero
Nil
a
anInstanceOf
True
False
*/

type numberCompare struct {
	actual interface{}
}

func (n *numberCompare) GreaterThan(bound float64) {

}

func (n *numberCompare) LessThan(bound float64) {

}

func (n *numberCompare) Within(low float64, high float64) {

}

type has struct {
	actual interface{}
	No     *listing
	Any    *listing
	All    *listing
	Length *length
}

func newHas(actual interface{}) *has {
	return &has{}
}

type length struct {
	*numberCompare
	Not *numberCompare
}

/*
Property
Key
DeepProperty
*/

type listing struct {
	actual interface{}
	state  bool
	out    bool
}

func (l *listing) Properties(names ...string) {
	params := append([]interface{}{}, l.actual, l.state, l.out)
	for _, name := range names {
		params = append(params, name)
	}
	expectations.check("has properties", true, params)
}

func (l *listing) Keys(names ...string) {
	params := append([]interface{}{}, l.actual, l.state, l.out)
	for _, name := range names {
		params = append(params, name)
	}
	expectations.check("has keys", true, params)
}

/*
Properties
Keys
*/
