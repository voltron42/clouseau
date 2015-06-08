package common_test

import (
	"."
	"../reckon"
	"../suiteshop"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	if suiteshop.Describe("Args", func(suite *suiteshop.Suite) {
		suite.Test("AddAll", func(log *suiteshop.Log) {
			args := &common.Args{}
			err := args.AddAll(reflect.ValueOf([]interface{}{1, 2, 3}))
			if err != nil {
				panic(err)
			}
			reckon.That(*args).Is.EqualTo([]interface{}{1, 2, 3})
		})
		suite.Test("Len & Add", func(log *suiteshop.Log) {
			args := &common.Args{1, 2}
			reckon.That(args.Len()).Is.EqualTo(2)
			args.Add(reflect.ValueOf(3))
			reckon.That(*args).Is.EqualTo([]interface{}{1, 2, 3})
			reckon.That(args.Len()).Is.EqualTo(3)
		})
	}).Post(func(message string) {
		list = append(list, message)
	}) {
		fmt.Println(strings.Join(list, "\n"))
	} else {
		t.Fatal(strings.Join(list, "\n"))
	}
}
