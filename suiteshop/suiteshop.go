package suiteshop

import (
	"../reckon"
	"errors"
	"fmt"
)

type Suite struct {
	label     string
	before    []dependency
	after     []dependency
	beforeAll []dependency
	afterAll  []dependency
	tests     []test
}

func (s *Suite) Before(fn func()) {
	s.before = append(s.before, dependency(fn))
}

func (s *Suite) BeforeAll(fn func()) {
	s.beforeAll = append(s.beforeAll, dependency(fn))
}

func (s *Suite) After(fn func()) {
	s.after = append(s.after, dependency(fn))
}

func (s *Suite) AfterAll(fn func()) {
	s.afterAll = append(s.afterAll, dependency(fn))
}

func (s *Suite) Test(label string, fn func()) {
	s.tests = append(s.tests, test{label, fn})
}

func Describe(label string, core func(suite *Suite)) *log {
	l := &log{&[]string{}}
	s := &Suite{
		label,
		[]dependency{},
		[]dependency{},
		[]dependency{},
		[]dependency{},
		[]test{},
	}
	core(s)
	for _, b := range s.beforeAll {
		b.runSafe(s.label)
	}
	for _, test := range s.tests {
		for _, before := range s.before {
			before.runSafe(s.label)
		}
		test.run(s.label)
		for _, after := range s.after {
			after.runSafe(s.label)
		}
	}
	for _, a := range s.afterAll {
		a.runSafe(s.label)
	}
	return l
}

type dependency func()

func (d dependency) runSafe(label string) {
	defer failPanic(label)
	d()
}

type test struct {
	label string
	fn    dependency
}

func (t test) run(label string) {
	t.fn.runSafe(label + " - " + t.label)
}

func failPanic(label string) {
	r := recover()
	if r != nil {
		switch r.(type) {
		case string:
			reckon.Fail(errors.New(fmt.Sprintf("%v -- %v", label, r)))
		case error:
			err, _ := r.(error)
			reckon.Fail(errors.New(fmt.Sprintf("%v -- %v", label, err.Error())))
		default:
			reckon.Fail(errors.New(fmt.Sprintf("%v -- unknown error: %v", label, r)))
		}
	}
}

type log struct {
	messages *[]string
}
