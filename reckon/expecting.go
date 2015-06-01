package reckon

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func SetLog(t *testing.T) {
	fail = func(err error) {
		t.Fatal(err)
	}
}

func SetLogToPanic() {
	fail = func(err error) {
		panic(err)
	}
}

var fail func(err error)

var expectations = expectationSet(map[string]expectation{
	"equal": expectation{
		Message:    "Items not equal:\n\tActual: %v\n\tExpected: %v\n",
		NotMessage: "Items equal:\n\tActual: %v\n\tExpected: %v\n",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"exists": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"matches": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"contains": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"is a": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"instance of": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"is zero": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"greater than": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"less than": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"within": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"has property": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"has key": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"has deep property": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
	"has properties": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			size := len(params)
			if size < 4 {
				return false, errors.New("")
			}
			actual := reflect.ValueOf(params[0])
			if actual.Kind() != reflect.Map {
				return false, errors.New("")
			}
			stateValue := reflect.ValueOf(params[1]).Interface()
			state, ok := stateValue.(bool)
			if !ok {
				return false, errors.New("")
			}
			outValue := reflect.ValueOf(params[2]).Interface()
			out, ok := outValue.(bool)
			if !ok {
				return false, errors.New("")
			}

			return false, nil
		},
	},
	"has keys": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(params ...interface{}) (bool, error) {
			return false, nil
		},
	},
})

type expectationSet map[string]expectation

func (e expectationSet) check(name string, state bool, params ...interface{}) {
	exp, ok := e[name]
	if ok {
		exp.check(state, params)
	}
}

type expectation struct {
	Message    string
	NotMessage string
	Condition  func(params ...interface{}) (bool, error)
}

func (e *expectation) checkTrue(params ...interface{}) {
	cond, err := e.Condition(params)
	if err != nil {
		fail(err)
	} else {
		if cond {
			fail(errors.New(fmt.Sprintf(e.Message, params...)))
		}
	}
}

func (e *expectation) checkFalse(params ...interface{}) {
	cond, err := e.Condition(params)
	if err != nil {
		fail(err)
	} else {
		if !cond {
			fail(errors.New(fmt.Sprintf(e.NotMessage, params...)))
		}
	}
}

func (e *expectation) check(state bool, params ...interface{}) {
	if state {
		e.checkTrue(params)
	} else {
		e.checkFalse(params)
	}
}
