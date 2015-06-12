package reckon

import (
	"../common"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var expectations = expectationSet(map[string]expectation{
	"equals": expectation{
		Message:    "Items not equal:\n\tActual: %v\n\tExpected: %v",
		NotMessage: "Items equal:\n\tActual: %v\n\tExpected: %v",
		Params:     []int{0, 1},
		Condition: func(args *common.Args) bool {
			return reflect.DeepEqual(args.Get(0), args.Get(1))
		},
	},
	"exists": expectation{
		Message:    "Does not exist",
		NotMessage: "Does exist",
		Condition: func(args *common.Args) bool {
			actual := args.Get(0).ValueOf()
			return !actual.IsNil() && actual.IsValid()
		},
	},
	"panics": expectation{
		Message:    "Has not panicked",
		NotMessage: "Has panicked",
		Condition: func(args *common.Args) bool {
			fn := args.Get(0).ValueOf()
			err := ""
			safeCall(fn, &err)
			if len(err) == 0 {
				return false
			}
			return true
		},
	},
	"panics with message": expectation{
		Message:    "Has panicked with incorrect message",
		NotMessage: "Has panicked with correct message",
		Condition: func(args *common.Args) bool {
			fn := args.Get(0).ValueOf()
			err := ""
			safeCall(fn, &err)
			if len(err) == 0 {
				panic("Has not panicked")
			}
			return reflect.DeepEqual(err, args.Get(1).Elem())
		},
	},
	"matches": expectation{
		Message:    "Does not match string",
		NotMessage: "Does match string",
		Condition: func(args *common.Args) bool {
			actual := args.Get(0).String()
			regex := args.Get(1).String()
			exp := regexp.MustCompile(regex)
			return exp.MatchString(actual)
		},
	},
	"contains": expectation{
		Message:    "Does not contain string",
		NotMessage: "Does contain string",
		Condition: func(args *common.Args) bool {
			actual := args.Get(0).String()
			needle := args.Get(1).String()
			return strings.Contains(actual, needle)
		},
	},
	"is a": expectation{
		Message:    "Is not kind of",
		NotMessage: "Is kind of",
		Condition: func(args *common.Args) bool {
			return args.Get(0).ValueOf().Kind() == args.Get(1).Elem()
		},
	},
	"instance of": expectation{
		Message:    "Is not an instance of",
		NotMessage: "Is an instance of",
		Condition: func(args *common.Args) bool {
			return reflect.DeepEqual(args.Get(0).TypeOf(), args.Get(1).Elem())
		},
	},
	"is zero": expectation{
		Message:    "Is not the zero value for type",
		NotMessage: "Is the zero value for type",
		Condition: func(args *common.Args) bool {
			myType := args.Get(0).TypeOf()
			myValue := args.Get(0).Elem()
			zeroValue := reflect.Zero(myType)
			return reflect.DeepEqual(myValue, zeroValue.Interface())
		},
	},
	"greater than": expectation{
		Message:    "Is not greater than",
		NotMessage: "Is greater than",
		Condition: func(args *common.Args) bool {
			myValue := args.Get(0).Float64()
			bound := args.Get(1).Float64()
			return myValue > bound
		},
	},
	"less than": expectation{
		Message:    "Is not less than",
		NotMessage: "Is less than",
		Condition: func(args *common.Args) bool {
			myValue := args.Get(0).Float64()
			bound := args.Get(1).Float64()
			return myValue < bound
		},
	},
	"within": expectation{
		Message:    "Is not within",
		NotMessage: "Is within",
		Condition: func(args *common.Args) bool {
			myValue := args.Get(0).Float64()
			low := args.Get(1).Float64()
			high := args.Get(2).Float64()
			return myValue > low && myValue < high
		},
	},
	"has keys": expectation{
		Message:    "Does not have keys",
		NotMessage: "Does have keys",
		Condition: func(args *common.Args) bool {
			if args.Len() < 3 {
				panic("Has no keys")
			}
			actual := args.Get(0).ValueOf()
			if actual.Kind() != reflect.Map {
				panic("Value is not a Map")
			}
			keys := []string{}
			for _, value := range actual.MapKeys() {
				keys = append(keys, value.String())
			}
			keyList := "," + strings.Join(keys, ",") + ","
			state := args.Get(1).Bool()
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
		Message:    "Does not have properties",
		NotMessage: "Does have properties",
		Condition: func(args *common.Args) bool {
			if args.Len() < 3 {
				panic("Has no properties")
			}
			actualType := args.Get(0).TypeOf()
			if actualType.Kind() != reflect.Struct {
				panic("Value is not a Struct")
			}
			count := actualType.NumField()
			keys := []string{}
			for index := 0; index < count; index++ {
				keys = append(keys, actualType.Field(index).Name)
			}
			keyList := "," + strings.Join(keys, ",") + ","
			state := args.Get(1).Bool()
			for _, param := range (*args)[2:] {
				key := fmt.Sprintf("%v", param)
				if state == strings.Contains(keyList, ","+key+",") {
					return state
				}
			}
			return !state
		},
	},
	"has property": expectation{
		Message:    "Does not have property",
		NotMessage: "Does have property",
		Condition: func(args *common.Args) bool {
			actual, key, values := getHasParams(args)
			value := getKey(actual, key)
			if !value.IsValid() {
				return false
			}
			if values.Len() == 0 {
				return true
			}
			matches := values.Some(func(item *common.Arg, index int) bool {
				return reflect.DeepEqual(value.Interface(), item.Elem())
			})
			return matches.Len() > 0
		},
	},
	"has deep property": expectation{
		Message:    "Does not have deep property",
		NotMessage: "Does have deep property",
		Condition: func(args *common.Args) bool {
			actual, key, values := getHasParams(args)
			keys := strings.Split(key.String(), ".")
			for _, k := range keys {
				actual = getKey(actual, reflect.ValueOf(k))
			}
			value := actual
			if !value.IsValid() {
				return false
			}
			if values.Len() == 0 {
				return true
			}
			matches := values.Some(func(item *common.Arg, index int) bool {
				return reflect.DeepEqual(value.Interface(), item.Elem())
			})
			return matches.Len() > 0
		},
	},
})

func getHasParams(args *common.Args) (reflect.Value, reflect.Value, *common.Args) {
	if args.Len() < 2 {
		panic("Missing key")
	}
	actual := args.Get(0).ValueOf()
	key := args.Get(1).ValueOf()
	values := args.Subset(2, -1)
	return actual, key, values
}

func getKey(obj, key reflect.Value) reflect.Value {
	for obj.Kind() == reflect.Interface {
		obj = obj.Elem()
	}
	if obj.Kind() == reflect.Map {
		fmt.Println("map")
		return obj.MapIndex(key)
	} else if obj.Kind() == reflect.Struct {
		fmt.Println("struct")
		return obj.FieldByName(key.String())
	} else if obj.Kind() == reflect.Array || obj.Kind() == reflect.Slice {
		fmt.Println("slice")
		index, err := strconv.Atoi(key.String())
		if err != nil {
			panic(err)
		}
		return obj.Index(index)
	} else {
		panic("Value is not a map, struct, slice, or array")
	}
}

func safeCall(fn reflect.Value, err interface{}) {
	defer func() {
		out := reflect.ValueOf(err).Elem()
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				out.Set(reflect.ValueOf(fmt.Sprintf("%v", r)))
			case error:
				er, ok := r.(error)
				if !ok {
					out.Set(reflect.ValueOf(fmt.Sprintf("unknown error: %v", r)))
				} else {
					out.Set(reflect.ValueOf(er.Error()))
				}
			default:
				out.Set(reflect.ValueOf(fmt.Sprintf("unknown error: %v", r)))
			}
		}
	}()
	fn.Call([]reflect.Value{})
}

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
	Params     []int
	Condition  func(args *common.Args) bool
}

func (e *expectation) checkBase(state bool, message string, params []interface{}) {
	args := common.Args(params)
	if state != e.Condition(&args) {
		panic(e.buildMessage(message, params))
	}
}

func (e *expectation) buildMessage(message string, params []interface{}) string {
	args := []interface{}{}
	for _, index := range e.Params {
		args = append(args, params[index])
	}
	return fmt.Sprintf(message, args...)
}

func (e *expectation) check(state bool, params []interface{}) {
	message := ""
	if state {
		message = e.Message
	} else {
		message = e.NotMessage
	}
	e.checkBase(state, message, params)
}
