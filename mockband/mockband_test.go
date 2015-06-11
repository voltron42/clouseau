package mockband_test

import (
	"."

	"../common"
	"../reckon"
	"../suiteshop"

	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	hasErrors := suiteshop.Describe("MockBand", func(suite *suiteshop.Suite) {
		suite.Describe("When", func(suite *suiteshop.Suite) {
			suite.Test("no return", func(log *suiteshop.Log) {
				strVal := "string value"
				numVal := 5678
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method1", strVal, numVal).Return()
				obj.Method1(strVal, numVal)
				reckon.That(mock.HasCalled("Method1", strVal, numVal).Once()).Is.True()
			})
			suite.Test("no params", func(log *suiteshop.Log) {
				strVal := "string value"
				numVal := 5678
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method2").Return(strVal, numVal, nil)
				str, num, err := obj.Method2()
				reckon.That(mock.HasCalled("Method2").Once()).Is.True()
				reckon.That(str).Is.EqualTo(strVal)
				reckon.That(num).Is.EqualTo(numVal)
				reckon.That(err).Is.Nil()
			})
			suite.Test("var args", func(log *suiteshop.Log) {
				strVal := "string value"
				numVal := "5678"
				retVal := "this is the return value"
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method3", strVal, numVal).Return(retVal)
				out := obj.Method3(strVal, numVal)
				reckon.That(mock.HasCalled("Method3", strVal, numVal).Once()).Is.True()
				reckon.That(out).Is.EqualTo(retVal)
			})
			suite.Test("panic", func(log *suiteshop.Log) {
				strVal := "string value"
				numVal := 5678
				panicMessage := "This is the panic message."
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method1", strVal, numVal).Panic(panicMessage)
				reckon.That(func() {
					obj.Method1(strVal, numVal)
				}).Will.PanicWith(panicMessage)
				reckon.That(mock.HasCalled("Method1", strVal, numVal).Once()).Is.True()
			})
			suite.Test("inject", func(log *suiteshop.Log) {
				strVal := "string value"
				result := ""
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method4", &result).Inject(strVal, 0)
				err := obj.Method4(&result)
				reckon.That(mock.HasCalled("Method4", &result).Once()).Is.True()
				reckon.That(err).Is.Nil()
				reckon.That(result).Is.EqualTo(strVal)
			})
			suite.Describe("Using Any", func(suite *suiteshop.Suite) {
				mock := NewMockObject()
				var obj Object = mock
				mock.When("Method1", common.Any(), common.Any()).Return()
				obj.Method1("This is the first string value.", 17)
				obj.Method1("This is the second string value.", 23)
				obj.Method1("This is the third string value.", 31)
				obj.Method1("This is the fourth string value.", 43)
				obj.Method1("This is the fifth string value.", 59)
				calls := mock.GetCalls("Method1", common.Any(), common.Any())
				params := calls.GetParams(0)
				reckon.That(params.String(0)).Is.EqualTo("This is the first string value.")
				reckon.That(params.Int(1)).Is.EqualTo(17)
				params = calls.GetParams(1)
				reckon.That(params.String(0)).Is.EqualTo("This is the second string value.")
				reckon.That(params.Int(1)).Is.EqualTo(23)
				params = calls.GetParams(2)
				reckon.That(params.String(0)).Is.EqualTo("This is the third string value.")
				reckon.That(params.Int(1)).Is.EqualTo(31)
				params = calls.GetParams(3)
				reckon.That(params.String(0)).Is.EqualTo("This is the fourth string value.")
				reckon.That(params.Int(1)).Is.EqualTo(43)
				params = calls.GetParams(4)
				reckon.That(params.String(0)).Is.EqualTo("This is the fifth string value.")
				reckon.That(params.Int(1)).Is.EqualTo(59)
			})
		})
	}).Post(func(message string) {
		list = append(list, message)
	})
	if hasErrors {
		t.Fatal(strings.Join(list, "\n"))
	} else {
		fmt.Println(strings.Join(list, "\n"))
	}
}

type Object interface {
	Method1(arg1 string, arg2 int)
	Method2() (string, int, error)
	Method3(params ...string) interface{}
	Method4(pointer interface{}) error
}

type MockObject struct {
	*mockband.Mock
}

func NewMockObject() *MockObject {
	return &MockObject{mockband.NewMock()}
}

func (o *MockObject) Method1(arg1 string, arg2 int) {
	o.Mock.Called("Method1", arg1, arg2)
}

func (o *MockObject) Method2() (string, int, error) {
	args := o.Mock.Called("Method2")
	return args.String(0), args.Int(1), args.Error(2)
}

func (o *MockObject) Method3(params ...string) interface{} {
	args := o.Mock.CalledVarArg("Method3", params)
	return args.Get(0)
}

func (o *MockObject) Method4(pointer interface{}) error {
	args := o.Mock.Called("Method4", pointer)
	return args.Error(0)
}
