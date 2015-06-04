package reckon

import (
	"../common"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var expectations = expectationSet(map[string]expectation{
	"equals": expectation{
		Message:    "Items not equal:\n\tActual: %v\n\tExpected: %v\n",
		NotMessage: "Items equal:\n\tActual: %v\n\tExpected: %v\n",
		Condition: func(args *common.Args) bool {
			fmt.Println("equal")
			return reflect.DeepEqual(args.Get(0), args.Get(1))
		},
	},
	"exists": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			actual := args.ValueOf(0)
			return !actual.IsNil() && actual.IsValid()
		},
	},
	"matches": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			actual := args.String(0)
			regex := args.String(1)
			exp := regexp.MustCompile(regex)
			return exp.MatchString(actual)
		},
	},
	"contains": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			actual := args.String(0)
			needle := args.String(1)
			return strings.Contains(actual, needle)
		},
	},
	"is a": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			return reflect.DeepEqual(args.ValueOf(0).Kind(), args.Get(1))
		},
	},
	"instance of": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			return reflect.DeepEqual(args.ValueOf(0).Type(), args.Get(1))
		},
	},
	"is zero": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			myType := args.TypeOf(0)
			myValue := args.ValueOf(0)
			return reflect.DeepEqual(myValue, reflect.Zero(myType))
		},
	},
	"greater than": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			myValue := args.Float64(0)
			bound := args.Float64(1)
			return myValue > bound
		},
	},
	"less than": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			myValue := args.Float64(0)
			bound := args.Float64(1)
			return myValue < bound
		},
	},
	"within": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			myValue := args.Float64(0)
			low := args.Float64(1)
			high := args.Float64(2)
			return myValue > low && myValue < high
		},
	},
	"has property": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			return false
		},
	},
	"has key": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			return false
		},
	},
	"has deep property": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			return false
		},
	},
	"has keys": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			if args.Len() < 3 {
				panic("")
			}
			actual := args.ValueOf(0)
			if actual.Kind() != reflect.Map {
				panic("")
			}
			keys := []string{}
			for _, value := range actual.MapKeys() {
				keys = append(keys, value.String())
			}
			keyList := "," + strings.Join(keys, ",") + ","
			state := args.Bool(1)
			for _, param := range (*args)[2:] {
				key := fmt.Sprintf("%v", param)
				if state == strings.Contains(keyList, ","+key+",") {
					return state
				}
			}
			return !state
		},
	},
	"has properties": expectation{
		Message:    "",
		NotMessage: "",
		Condition: func(args *common.Args) bool {
			if args.Len() < 3 {
				panic("")
			}
			actualType := args.TypeOf(0)
			if actualType.Kind() != reflect.Struct {
				panic("")
			}
			count := actualType.NumField()
			keys := []string{}
			for index := 0; index < count; index++ {
				keys = append(keys, actualType.Field(index).Name)
			}
			keyList := "," + strings.Join(keys, ",") + ","
			state := args.Bool(1)
			for _, param := range (*args)[2:] {
				key := fmt.Sprintf("%v", param)
				if state == strings.Contains(keyList, ","+key+",") {
					return state
				}
			}
			return !state
		},
	},
})

type expectationSet map[string]expectation

func (e expectationSet) check(name string, state bool, params ...interface{}) {
	fmt.Println("checking " + name)
	exp, ok := e[name]
	if ok {
		fmt.Println("running check")
		exp.check(state, params)
	} else {
		fmt.Println("falling down")
	}
}

type expectation struct {
	Message    string
	NotMessage string
	Condition  func(args *common.Args) bool
}

func (e *expectation) checkBase(state bool, message string, params []interface{}) {
	args := common.Args(params)
	if state != e.Condition(&args) {
		panic(fmt.Sprintf(message, params...))
	}
}

func (e *expectation) check(state bool, params []interface{}) {
	message := ""
	if state {
		message = e.Message
	} else {
		message = e.NotMessage
	}
	fmt.Println(message)
	e.checkBase(state, message, params)
}
