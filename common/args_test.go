package common_test

import (
	"."
	"../reckon"
	"../suiteshop"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	hasErrors := suiteshop.Describe("Args", func(suite *suiteshop.Suite) {
		suite.Test("AddAll", func(log *suiteshop.Log) {
			args := &common.Args{}
			err := args.AddAll(reflect.ValueOf([]interface{}{1, 2, 3}))
			if err != nil {
				panic(err)
			}
			reckon.That(*args).Is.EqualTo(common.Args([]interface{}{1, 2, 3}))
		})
		suite.Test("Len & Add", func(log *suiteshop.Log) {
			args := &common.Args{1, 2}
			reckon.That(args.Len()).Is.EqualTo(2)
			args.Add(reflect.ValueOf(3))
			reckon.That(*args).Is.EqualTo(common.Args([]interface{}{1, 2, 3}))
			reckon.That(args.Len()).Is.EqualTo(3)
		})
		suite.Test("Get", func(log *suiteshop.Log) {
			obj := struct {
				a string
				b int
				c bool
			}{
				a: "goodbye",
				b: 17,
				c: false,
			}
			array := []int{7, 8, 9}
			myMap := map[string]int{
				"x": 10,
				"y": 11,
				"z": 12,
			}
			args := &common.Args{5, "hello", 1.256, true, obj, array, myMap}
			reckon.That(args.Get(0)).Is.EqualTo(5)
			reckon.That(args.Get(1)).Is.EqualTo("hello")
			reckon.That(args.Get(2)).Is.EqualTo(1.256)
			reckon.That(args.Get(3)).Is.EqualTo(true)
			reckon.That(args.Get(4)).Is.EqualTo(obj)
			reckon.That(args.Get(5)).Is.EqualTo(array)
			reckon.That(args.Get(6)).Is.EqualTo(myMap)
			reckon.That(args.Get(7)).Is.Nil()
		})
		suite.Describe("Match", func(suite *suiteshop.Suite) {
			suite.Test("any success", func(log *suiteshop.Log) {
				args1 := &common.Args{1, 2, 3}
				args2 := &common.Args{common.Any(), common.Any(), common.Any()}
				reckon.That(args2.Matches(args1)).Is.True()
				reckon.That(args1.Matches(args2)).Is.False()
			})
			suite.Test("some success", func(log *suiteshop.Log) {
				args1 := &common.Args{1, 2, 3}
				args2 := &common.Args{1, 2, common.Any()}
				args3 := &common.Args{1, common.Any(), 3}
				args4 := &common.Args{common.Any(), 2, 3}
				args5 := &common.Args{4, 2, 3}
				reckon.That(args2.Matches(args1)).Is.True()
				reckon.That(args3.Matches(args1)).Is.True()
				reckon.That(args4.Matches(args1)).Is.True()
				reckon.That(args5.Matches(args1)).Is.False()
			})
		})
		suite.Describe("Casting", func(suite *suiteshop.Suite) {
			suite.Test("Numbers", func(log *suiteshop.Log) {
				args := &common.Args{-71215.23546873}
				reckon.That(args.Get(0)).Is.EqualTo(-71215.23546873)
				reckon.That(args.String(0)).Is.EqualTo("-71215.23546873")
				reckon.That(args.Byte(0)).Is.EqualTo(byte(209))
				reckon.That(args.Float32(0)).Is.EqualTo(float32(-71215.234))
				reckon.That(args.Float64(0)).Is.EqualTo(-71215.23546873)
				reckon.That(args.Int(0)).Is.EqualTo(-71215)
				reckon.That(args.Int8(0)).Is.EqualTo(int8(-47))
				reckon.That(args.Int16(0)).Is.EqualTo(int16(-5679))
				reckon.That(args.Int32(0)).Is.EqualTo(int32(-71215))
				reckon.That(args.Int64(0)).Is.EqualTo(int64(-71215))
				reckon.That(args.UInt(0)).Is.EqualTo(uint(18446744073709480401))
				reckon.That(args.UInt16(0)).Is.EqualTo(uint16(59857))
				reckon.That(args.UInt32(0)).Is.EqualTo(uint32(4294896081))
				reckon.That(args.UInt64(0)).Is.EqualTo(uint64(18446744073709480401))
				reckon.That(args.UIntPtr(0)).Is.EqualTo(uintptr(18446744073709480401))
			})
			suite.Test("Strings", func(log *suiteshop.Log) {
				obj := struct {
					a string
					b int
					c bool
				}{
					a: "goodbye",
					b: 17,
					c: false,
				}
				array := []int{7, 8, 9}
				myMap := map[string]int{
					"x": 10,
					"y": 11,
					"z": 12,
				}
				args := &common.Args{5, "hello", 1.256, true, obj, array, myMap}
				reckon.That(args.String(0)).Is.EqualTo("5")
				reckon.That(args.String(1)).Is.EqualTo("hello")
				reckon.That(args.String(2)).Is.EqualTo("1.256")
				reckon.That(args.String(3)).Is.EqualTo("true")
				reckon.That(args.String(4)).Is.EqualTo("{goodbye 17 false}")
				reckon.That(args.String(5)).Is.EqualTo("[7 8 9]")
				reckon.That(args.String(6)).Is.EqualTo("map[x:10 y:11 z:12]")
			})
			suite.Test("Boolean", func(log *suiteshop.Log) {
				args := &common.Args{true, false}
				reckon.That(args.Bool(0)).Is.True()
				reckon.That(args.Bool(1)).Is.False()
			})
			suite.Test("Error", func(log *suiteshop.Log) {
				args := &common.Args{errors.New("hi there")}
				reckon.That(args.Error(0).Error()).Is.EqualTo("hi there")
			})
			suite.Test("Bytes", func(log *suiteshop.Log) {
				args := &common.Args{[]byte("hi there")}
				reckon.That(args.Bytes(0)).Is.EqualTo([]byte{104, 105, 32, 116, 104, 101, 114, 101})
			})
			suite.Test("Strings", func(log *suiteshop.Log) {
				args := &common.Args{strings.Split("The quick brown fox jumped over the lazy dog.", " ")}
				list := args.Strings(0)
				reckon.That(list[0]).Is.EqualTo("The")
				reckon.That(list[1]).Is.EqualTo("quick")
				reckon.That(list[2]).Is.EqualTo("brown")
				reckon.That(list[3]).Is.EqualTo("fox")
				reckon.That(list[4]).Is.EqualTo("jumped")
				reckon.That(list[5]).Is.EqualTo("over")
				reckon.That(list[6]).Is.EqualTo("the")
				reckon.That(list[7]).Is.EqualTo("lazy")
				reckon.That(list[8]).Is.EqualTo("dog.")
			})
			suite.Test("Fail", func(log *suiteshop.Log) {
				args := &common.Args{}
				reckon.That(func() {
					args.Bytes(0)
				}).Will.PanicWith("Cannot be cast to byte array.")
				reckon.That(func() {
					args.Error(0)
				}).Will.PanicWith("Cannot be cast to error.")
				reckon.That(func() {
					args.Strings(0)
				}).Will.PanicWith("Cannot be cast to string array.")
			})
		})
	}).Post(func(message string) {
		list = append(list, message)
	})

	if hasErrors {
		fmt.Println(strings.Join(list, "\n"))
	} else {
		t.Fatal(strings.Join(list, "\n"))
	}
}
