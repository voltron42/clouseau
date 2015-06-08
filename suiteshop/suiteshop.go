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
	for _, b := range s.beforeAll {
		b.runSafe(label, log)
	}
	for _, test := range s.tests {
		for _, before := range s.before {
			before.runSafe(label, log)
		}
		test.run(label, log)
		for _, after := range s.after {
			after.runSafe(label, log)
		}
	}
	for _, suite := range s.suites {
		suite.run(log, allLabels...)
	}
	for _, a := range s.afterAll {
		a.runSafe(label, log)
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

func (d dependency) runSafe(label string, log *Log) {
	defer failPanic(label, log)
	d(log)
}

type test struct {
	label string
	fn    dependency
}

func (t test) run(label string, log *Log) {
	testName := label + " - " + t.label
	log.append(fmt.Sprintf("\x1b[42;37;1m%v\x1b[0m", testName))
	t.fn.runSafe(testName, log)
	log.append("")
}

func failPanic(label string, log *Log) {
	r := recover()
	if r != nil {
		switch r.(type) {
		case string:
			log.append(fmt.Sprintf("\x1b[41;37;1m%v -- %v\x1b[0m", label, r))
		case error:
			err, _ := r.(error)
			log.append(fmt.Sprintf("\x1b[41;37;1m%v -- %v\x1b[0m", label, err.Error()))
		default:
			log.append(fmt.Sprintf("\x1b[41;37;1m%v -- unknown error: %v\x1b[0m", label, r))
		}
		stack := string(debug.Stack())
		stack = strings.Join(strings.Split(stack, "\n"), "\n\n")
		stack = strings.Join(strings.Split(stack, "\n\n\t"), "\n\t")
		log.append(fmt.Sprintf("\n\x1b[41;37;1m%v\x1b[0m", stack))
	}
}

type Log struct {
	hasErrors bool
	messages  []string
	info      []string
}

func newLog() *Log {
	return &Log{false, []string{}, []string{}}
}

func (l *Log) Info(message string) {
	l.info = append(l.info, message)
}

func (l *Log) append(message string) {
	l.hasErrors = true
	l.messages = append(l.messages, message)
}

func (log *Log) Post(fn func(message string)) bool {
	fn("")
	for _, message := range log.info {
		fn(message)
	}
	fn("")
	for _, message := range log.messages {
		fn(message)
	}
	fn("")
	return log.hasErrors
}
