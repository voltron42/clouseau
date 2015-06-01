package reckon

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var T *testing.T

var Fail = func(err error) {
	T.Fatal(err)
}

func That(actual interface{}) *reckoning {
	return &reckoning{actual}
}

type reckoning struct {
	actual interface{}
	Is is
	Has has
}

type is struct {
  	actual interface{}
}

func (i *is) EqualTo(expected interface{}) {
	if !reflect.DeepEqual(i.actual, expected) {
		Fail(errors.New(fmt.Sprintf("Reckoning Failed:\n\tActual: %v\n\tExpected: %v\n", i.actual, expected)))
	}
}



