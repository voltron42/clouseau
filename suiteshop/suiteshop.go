package suiteshop

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type Suite struct {
	label     string
	before    []dependency
	after     []dependency
	beforeAll []dependency
	afterAll  []dependency
	tests     []test
	suites    []Suite
}

func newSuite(label string) *Suite {
	return &Suite{
		label,
		[]dependency{},
		[]dependency{},
		[]dependency{},
		[]dependency{},
		[]test{},
		[]Suite{},
	}
}

func (s *Suite) Before(fn func(log *Log)) {
	s.before = append(s.before, dependency(fn))
}

func (s *Suite) BeforeAll(fn func(log *Log)) {
	s.beforeAll = append(s.beforeAll, dependency(fn))
}

func (s *Suite) After(fn func(log *Log)) {
	s.after = append(s.after, dependency(fn))
}

func (s *Suite) AfterAll(fn func(log *Log)) {
	s.afterAll = append(s.afterAll, dependency(fn))
}

func (s *Suite) Test(label string, fn func(log *Log)) {
	s.tests = append(s.tests, test{label, fn})
}

func (s *Suite) Describe(label string, core func(suite *Suite)) {
	s2 := newSuite(label)
	core(s2)
	s.suites = append(s.suites, *s2)
}

func (s *Suite) run(log *Log, labels ...string) {
	allLabels := append(labels, s.label)
	label := strings.Join(allLabels, " - ")
	test := &testRun{label, nil}
	for _, b := range s.beforeAll {
		b.runSafe(test, log)
		if test.exception != nil {
			break
		}
	}
	if test.exception == nil {
		for _, t := range s.tests {
			for _, before := range s.before {
				before.runSafe(test, log)
				if test.exception != nil {
					break
				}
			}
			if test.exception == nil {
				run := t.run(label, log)
				log.append(run)
				if run.exception != nil {
					continue
				} else {
					for _, after := range s.after {
						after.runSafe(test, log)
						if test.exception != nil {
							break
						}
					}
				}
			}
		}
		for _, suite := range s.suites {
			suite.run(log, allLabels...)
		}
		for _, a := range s.afterAll {
			a.runSafe(test, log)
			if test.exception != nil {
				break
			}
		}
	}
	if test.exception != nil {
		log.append(test)
	}
}

func Describe(label string, core func(suite *Suite)) *Log {
	s := newSuite(label)
	core(s)
	log := newLog()
	s.run(log)
	return log
}

type dependency func(log *Log)

func (d dependency) runSafe(test *testRun, log *Log) {
	defer failPanic(test, log)
	d(log)
}

type test struct {
	label string
	fn    dependency
}

func (t test) run(label string, log *Log) *testRun {
	testName := label + " - " + t.label
	test := &testRun{testName, nil}
	t.fn.runSafe(test, log)
	return test
}

func failPanic(test *testRun, log *Log) {
	if r := recover(); r != nil {
		message := ""
		switch r.(type) {
		case string:
			message = fmt.Sprintf("%v -- %v", test.label, r)
		case error:
			err, _ := r.(error)
			message = (fmt.Sprintf("%v -- %v", test.label, err.Error()))
		default:
			message = (fmt.Sprintf("%v -- unknown error: %v", test.label, r))
		}
		test.setException(message, string(debug.Stack()))
	}
}

type stackStep struct {
	location string
	call     string
}

func (s *stackStep) String() string {
	return fmt.Sprintf("%v\n\t%v", s.location, s.call)
}

type exception struct {
	message string
	stack   []stackStep
}

func (e *exception) String() string {
	list := []string{e.message}
	for _, item := range e.stack {
		list = append(list, item.String())
	}
	return fmt.Sprintf("\x1b[41;37;1m%v\x1b[0m\n", strings.Join(list, "\n\n"))
}

type testRun struct {
	label     string
	exception *exception
}

func (t *testRun) setException(message, stack string) {
	stack = strings.Join(strings.Split(stack, "\n\t"), "\t")
	list := strings.Split(stack, "\n")
	steps := []stackStep{}
	for _, step := range list {
		fields := strings.Split(step, "\t")
		if len(fields) >= 2 {
			steps = append(steps, stackStep{fields[0], fields[1]})
		}
	}
	t.exception = &exception{message, steps}
}

func (t *testRun) String() string {
	e := ""
	if t.exception != nil {
		e = t.exception.String()
	}
	return fmt.Sprintf("\x1b[42;37;1m%v\x1b[0m\n%v\n", t.label, e)
}

type Log struct {
	tests    []testRun
	messages []string
}

func (l *Log) hasErrors() bool {
	for _, test := range l.tests {
		if test.exception != nil {
			fmt.Println(test.exception)
			return true
		}
	}
	return false
}

func newLog() *Log {
	return &Log{[]testRun{}, []string{}}
}

func (l *Log) Info(message string) {
	l.messages = append(l.messages, message)
}

func (l *Log) append(test *testRun) {
	l.tests = append(l.tests, *test)
}

func (log *Log) Post(fn func(message string)) bool {
	fn("")
	for _, message := range log.messages {
		fn(message)
	}
	fn("")
	for _, test := range log.tests {
		fn(test.String())
	}
	fn("")
	return log.hasErrors()
}
