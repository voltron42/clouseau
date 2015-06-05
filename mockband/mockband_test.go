package mockband_test

import (
	"."

	"../reckon"
	"../suiteshop"

	"fmt"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	suiteshop.Describe("MockBand", func(suite *suiteshop.Suite) {
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
		})
	}).Post(func(message string) {
		fmt.Println(message)
		list = append(list, message)
	})
}

type Object interface {
	Method1(arg1 string, arg2 int)
	Method2() (string, int, error)
	Method3(params ...string) interface{}
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
