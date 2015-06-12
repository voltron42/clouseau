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
		suite.Test("PopLast", func(log *suiteshop.Log) {
			args := &common.Args{1, 2}
			last := args.PopLast()
			reckon.That(last.Elem()).Is.EqualTo(2)
			reckon.That(args.Len()).Is.EqualTo(1)
			reckon.That(args.Get(0).Elem()).Is.EqualTo(1)
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
			reckon.That(args.Get(0).Elem()).Is.EqualTo(5)
			reckon.That(args.Get(1).Elem()).Is.EqualTo("hello")
			reckon.That(args.Get(2).Elem()).Is.EqualTo(1.256)
			reckon.That(args.Get(3).Elem()).Is.EqualTo(true)
			reckon.That(args.Get(4).Elem()).Is.EqualTo(obj)
			reckon.That(args.Get(5).Elem()).Is.EqualTo(array)
			reckon.That(args.Get(6).Elem()).Is.EqualTo(myMap)
			reckon.That(args.Get(7).Elem()).Is.Nil()
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
				arg := args.Get(0)
				reckon.That(arg.Elem()).Is.EqualTo(-71215.23546873)
				reckon.That(arg.String()).Is.EqualTo("-71215.23546873")
				reckon.That(arg.Byte()).Is.EqualTo(byte(209))
				reckon.That(arg.Float32()).Is.EqualTo(float32(-71215.234))
				reckon.That(arg.Float64()).Is.EqualTo(-71215.23546873)
				reckon.That(arg.Int()).Is.EqualTo(-71215)
				reckon.That(arg.Int8()).Is.EqualTo(int8(-47))
				reckon.That(arg.Int16()).Is.EqualTo(int16(-5679))
				reckon.That(arg.Int32()).Is.EqualTo(int32(-71215))
				reckon.That(arg.Int64()).Is.EqualTo(int64(-71215))
				reckon.That(arg.UInt()).Is.EqualTo(uint(18446744073709480401))
				reckon.That(arg.UInt16()).Is.EqualTo(uint16(59857))
				reckon.That(arg.UInt32()).Is.EqualTo(uint32(4294896081))
				reckon.That(arg.UInt64()).Is.EqualTo(uint64(18446744073709480401))
				reckon.That(arg.UIntPtr()).Is.EqualTo(uintptr(18446744073709480401))
			})
			suite.Test("String", func(log *suiteshop.Log) {
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
				reckon.That(args.Get(0).String()).Is.EqualTo("5")
				reckon.That(args.Get(1).String()).Is.EqualTo("hello")
				reckon.That(args.Get(2).String()).Is.EqualTo("1.256")
				reckon.That(args.Get(3).String()).Is.EqualTo("true")
				reckon.That(args.Get(4).String()).Is.EqualTo("{goodbye 17 false}")
				reckon.That(args.Get(5).String()).Is.EqualTo("[7 8 9]")
				reckon.That(args.Get(6).String()).Is.EqualTo("map[x:10 y:11 z:12]")
			})
			suite.Test("Boolean", func(log *suiteshop.Log) {
				args := &common.Args{true, false}
				reckon.That(args.Get(0).Bool()).Is.True()
				reckon.That(args.Get(1).Bool()).Is.False()
			})
			suite.Test("Error", func(log *suiteshop.Log) {
				args := &common.Args{errors.New("hi there")}
				reckon.That(args.Get(0).Error().Error()).Is.EqualTo("hi there")
			})
			suite.Test("Bytes", func(log *suiteshop.Log) {
				args := &common.Args{[]byte("hi there")}
				reckon.That(args.Get(0).Bytes()).Is.EqualTo([]byte{104, 105, 32, 116, 104, 101, 114, 101})
			})
			suite.Test("Strings", func(log *suiteshop.Log) {
				args := &common.Args{strings.Split("The quick brown fox jumped over the lazy dog.", " ")}
				list := args.Get(0).Strings()
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
				arg := args.Get(0)
				reckon.That(arg.Bytes()).Is.Zero()
				reckon.That(arg.Error()).Is.Nil()
				reckon.That(arg.Strings()).Is.Zero()
				reckon.That(func() {
					arg.Bool()
				}).Will.Panic()
			})
		})
		suite.Test("Some", func(log *suiteshop.Log) {
			args := &common.Args{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			getEvens := func(item *common.Arg, index int) bool {
				return item.Int()%2 == 0
			}
			evens := &common.Args{2, 4, 6, 8, 10}
			reckon.That(args.Some(getEvens)).Is.EqualTo(evens)
			getBottomHalf := func(item *common.Arg, index int) bool {
				return item.Int() <= 5
			}
			bottomHalf := &common.Args{1, 2, 3, 4, 5}
			reckon.That(args.Some(getBottomHalf)).Is.EqualTo(bottomHalf)
		})
		suite.Test("Map", func(log *suiteshop.Log) {
			args := &common.Args{1, 2, 3, 4, 5}
			double := func(item *common.Arg, index int) *common.Arg {
				return common.NewArg(item.Int() * 2)
			}
			doubles := &common.Args{2, 4, 6, 8, 10}
			reckon.That(args.Map(double)).Is.EqualTo(doubles)
		})
		suite.Test("Each", func(log *suiteshop.Log) {
			args := &common.Args{1, 2, 3, 4, 5}
			sum := 0
			calcSum := func(item *common.Arg, index int) {
				sum += item.Int()
			}
			args.Each(calcSum)
			reckon.That(sum).Is.EqualTo(15)
		})
		suite.Test("Subset", func(log *suiteshop.Log) {
			args := &common.Args{1, 2, 3, 4, 5}
			reckon.That(args.Subset(0, 2)).Is.EqualTo(&common.Args{1, 2})
			reckon.That(args.Subset(3, -1)).Is.EqualTo(&common.Args{4, 5})
			reckon.That(args.Subset(2, 3)).Is.EqualTo(&common.Args{3})
			reckon.That(args.Subset(2, 2)).Is.EqualTo(&common.Args{})
			reckon.That(func() {
				args.Subset(2, 1)
			}).Will.PanicWith("runtime error: slice bounds out of range")
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
