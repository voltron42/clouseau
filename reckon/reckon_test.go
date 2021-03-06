package reckon_test

import (
	"."
	"../suiteshop"

	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	fn := func(message string) {
		list = append(list, message)
	}
	if suiteshop.Describe("Reckon", func(suite *suiteshop.Suite) {
		suite.Describe("Panic", func(suite *suiteshop.Suite) {
			suite.Test("will panic", func(log *suiteshop.Log) {
				reckon.That(func() {
					panic("panicking")
				}).Will.Panic()
			})
			suite.Test("will panic with wrong message", func(log *suiteshop.Log) {
				defer func() {
					if r := recover(); r == nil {
						panic("should have panicked with message")
					} else {
						if r != "Has panicked with incorrect message" {
							panic(fmt.Sprintf("Has not panicked with incorrect message: %v", r))
						}
					}
				}()
				reckon.That(func() {
					panic("panicking")
				}).Will.PanicWith("should panic")
			})
		})
		suite.Describe("Not Panic", func(suite *suiteshop.Suite) {
			suite.Test("will panic", func(log *suiteshop.Log) {
				defer func() {
					if r := recover(); r == nil {
						panic("should have panicked with message")
					} else {
						if r != "Has panicked" {
							panic(fmt.Sprintf("Has not panicked with incorrect message: %v", r))
						}
					}
				}()
				reckon.That(func() {
					panic("panicking")
				}).Will.Not.Panic()
			})
		})
		suite.Describe("Equals", func(suite *suiteshop.Suite) {
			suite.Test("is equal", func(log *suiteshop.Log) {
				reckon.That(5).Is.EqualTo(5)
			})
			suite.Test("not equal", func(log *suiteshop.Log) {
				reckon.That(func() {
					reckon.That(4).Is.EqualTo(5)
				}).Will.PanicWith("Items not equal:\n\tActual: 4\n\tExpected: 5")
			})
		})
		suite.Describe("Not Equal", func(suite *suiteshop.Suite) {
			suite.Test("is equal", func(log *suiteshop.Log) {
				reckon.That(func() {
					reckon.That(5).Is.Not.EqualTo(5)
				}).Will.PanicWith("Items equal:\n\tActual: 5\n\tExpected: 5")
			})
			suite.Test("not equal", func(log *suiteshop.Log) {
				reckon.That(4).Is.Not.EqualTo(5)
			})
		})
		suite.Describe("Is [Something]", func(suite *suiteshop.Suite) {
			suite.Describe("Nil", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(nil).Is.Nil()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("hi").Is.Nil()
					}).Will.Panic()
				})
			})
			suite.Describe("False", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(false).Is.False()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("hi").Is.False()
					}).Will.Panic()
				})
			})
			suite.Describe("True", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(true).Is.True()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("hi").Is.True()
					}).Will.Panic()
				})
			})
			suite.Describe("Zero", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(0.0).Is.Zero()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.Zero()
					}).Will.Panic()
				})
			})
			suite.Describe("A [Kind]", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(0.0).Is.A(reflect.Float64)
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.A(reflect.String)
					}).Will.PanicWith("Is not kind of")
				})
			})
			suite.Describe("A Instance Of [Type]", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(0.0).Is.AnInstanceOf(reflect.TypeOf(1.0))
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.AnInstanceOf(reflect.TypeOf(""))
					}).Will.PanicWith("Is not an instance of")
				})
			})
		})
		suite.Describe("Is Not [Something]", func(suite *suiteshop.Suite) {
			suite.Describe("Nil", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That("nil").Is.Not.Nil()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(nil).Is.Not.Nil()
					}).Will.Panic()
				})
			})
			suite.Describe("False", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That("false").Is.Not.False()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(false).Is.Not.False()
					}).Will.Panic()
				})
			})
			suite.Describe("True", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That("true").Is.Not.True()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(true).Is.Not.True()
					}).Will.Panic()
				})
			})
			suite.Describe("Zero", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(5).Is.Not.Zero()
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("").Is.Not.Zero()
					}).Will.Panic()
				})
			})
			suite.Describe("A [Kind]", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(5).Is.Not.A(reflect.Float64)
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("5").Is.Not.A(reflect.String)
					}).Will.PanicWith("Is kind of")
				})
			})
			suite.Describe("A Instance Of [Type]", func(suite *suiteshop.Suite) {
				suite.Test("success", func(log *suiteshop.Log) {
					reckon.That(true).Is.Not.AnInstanceOf(reflect.TypeOf(1.0))
				})
				suite.Test("fail", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That("5").Is.Not.AnInstanceOf(reflect.TypeOf(""))
					}).Will.PanicWith("Is an instance of")
				})
			})
		})
		suite.Describe("Numeric", func(suite *suiteshop.Suite) {
			suite.Describe("Greater Than", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(5).Is.GreaterThan(4)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(4).Is.GreaterThan(5)
					}).Will.Panic()
				})
			})
			suite.Describe("Not Greater Than", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(4).Is.Not.GreaterThan(4)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.Not.GreaterThan(4)
					}).Will.Panic()
				})
			})
			suite.Describe("Less Than", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(4).Is.LessThan(5)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.LessThan(4)
					}).Will.Panic()
				})
			})
			suite.Describe("Not Less Than", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(5).Is.Not.LessThan(5)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.Not.LessThan(6)
					}).Will.Panic()
				})
			})
			suite.Describe("Within", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(4).Is.Within(3, 5)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.Within(3, 4)
					}).Will.Panic()
				})
			})
			suite.Describe("Not Within", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(4).Is.Not.Within(3, 3)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(5).Is.Not.Within(3, 6)
					}).Will.Panic()
				})
			})
		})
		suite.Describe("Has", func(suite *suiteshop.Suite) {
			suite.Describe("Any", func(suite *suiteshop.Suite) {
				suite.Describe("Keys", func(suite *suiteshop.Suite) {
					suite.Test("does", func(log *suiteshop.Log) {
						reckon.That(map[string]int{
							"a": 1,
							"b": 2,
							"c": 3,
						}).Has.Any.Keys("c", "d")
					})
					suite.Test("does not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That(map[string]int{
								"a": 1,
								"b": 2,
								"c": 3,
							}).Has.Any.Keys("d")
						}).Will.PanicWith("Does not have keys")
					})
				})
				suite.Describe("Properties", func(suite *suiteshop.Suite) {
					suite.Test("does", func(log *suiteshop.Log) {
						reckon.That(struct {
							a int
							b int
							c int
						}{
							a: 1,
							b: 2,
							c: 3,
						}).Has.Any.Properties("c", "d")
					})
					suite.Test("does not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That(struct {
								a int
								b int
								c int
							}{
								a: 1,
								b: 2,
								c: 3,
							}).Has.Any.Properties("d")
						}).Will.PanicWith("Does not have properties")
					})
				})
			})
			suite.Describe("All", func(suite *suiteshop.Suite) {
				suite.Describe("Keys", func(suite *suiteshop.Suite) {
					suite.Test("does", func(log *suiteshop.Log) {
						reckon.That(map[string]int{
							"a": 1,
							"b": 2,
							"c": 3,
						}).Has.All.Keys("c", "b")
					})
					suite.Test("does not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That(map[string]int{
								"a": 1,
								"b": 2,
								"c": 3,
							}).Has.All.Keys("c", "d")
						}).Will.PanicWith("Does not have keys")
					})
				})
				suite.Describe("Properties", func(suite *suiteshop.Suite) {
					suite.Test("does", func(log *suiteshop.Log) {
						reckon.That(struct {
							a int
							b int
							c int
						}{
							a: 1,
							b: 2,
							c: 3,
						}).Has.All.Properties("c", "b")
					})
					suite.Test("does not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That(struct {
								a int
								b int
								c int
							}{
								a: 1,
								b: 2,
								c: 3,
							}).Has.All.Properties("c", "d")
						}).Will.PanicWith("Does not have properties")
					})
				})
			})
			suite.Describe("Length", func(suite *suiteshop.Suite) {
				suite.Describe("Greater Than", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2}).Has.Length.GreaterThan(1)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2}).Has.Length.GreaterThan(5)
						}).Will.Panic()
					})
				})
				suite.Describe("Not Greater Than", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2, 3, 4}).Has.Length.Not.GreaterThan(4)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2, 3}).Has.Length.Not.GreaterThan(2)
						}).Will.Panic()
					})
				})
				suite.Describe("Less Than", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2}).Has.Length.LessThan(3)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2, 3}).Has.Length.LessThan(2)
						}).Will.Panic()
					})
				})
				suite.Describe("Not Less Than", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2}).Has.Length.Not.LessThan(2)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2}).Has.Length.Not.LessThan(6)
						}).Will.Panic()
					})
				})
				suite.Describe("Within", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2, 3, 4}).Has.Length.Within(3, 5)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2}).Has.Length.Within(3, 4)
						}).Will.Panic()
					})
				})
				suite.Describe("Not Within", func(suite *suiteshop.Suite) {
					suite.Test("is", func(log *suiteshop.Log) {
						reckon.That([]int{1, 2}).Has.Length.Not.Within(3, 3)
					})
					suite.Test("is not", func(log *suiteshop.Log) {
						reckon.That(func() {
							reckon.That([]int{1, 2, 3, 4}).Has.Length.Not.Within(3, 6)
						}).Will.Panic()
					})
				})
			})
			suite.Describe("Property", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(map[int]string{3: "a"}).Has.Property(3)
					reckon.That(map[int]string{3: "a"}).Has.Property(3, "a")
					reckon.That(map[string]int{"a": 3}).Has.Property("a", 3)
					reckon.That(map[string]int{"a": 3}).Has.Property("a")
					reckon.That(struct{ A int }{A: 3}).Has.Property("A")
					reckon.That(struct{ A int }{A: 3}).Has.Property("A", 3)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(map[int]string{3: "a"}).Has.Property(2)
					}).Will.PanicWith("Does not have property")
					reckon.That(func() {
						reckon.That(map[int]string{3: "a"}).Has.Property(3, "b")
					}).Will.PanicWith("Does not have property")
					reckon.That(func() {
						reckon.That(map[string]int{"a": 3}).Has.Property("b")
					}).Will.PanicWith("Does not have property")
					reckon.That(func() {
						reckon.That(map[string]int{"a": 3}).Has.Property("a", 2)
					}).Will.PanicWith("Does not have property")
					reckon.That(func() {
						reckon.That(struct{ A int }{A: 3}).Has.Property("B")
					}).Will.PanicWith("Does not have property")
					reckon.That(func() {
						reckon.That(struct{ A int }{A: 3}).Has.Property("A", 2)
					}).Will.PanicWith("Does not have property")
				})
			})
			suite.Describe("Not Property", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					reckon.That(map[int]string{3: "a"}).Has.No.Property(2)
					reckon.That(map[int]string{3: "a"}).Has.No.Property(3, "b")
					reckon.That(map[string]int{"a": 3}).Has.No.Property("b")
					reckon.That(map[string]int{"a": 3}).Has.No.Property("a", 2)
					reckon.That(struct{ A int }{A: 3}).Has.No.Property("B")
					reckon.That(struct{ A int }{A: 3}).Has.No.Property("A", 2)
				})
				suite.Test("is not", func(log *suiteshop.Log) {
					reckon.That(func() {
						reckon.That(map[int]string{3: "a"}).Has.No.Property(3)
					}).Will.PanicWith("Does have property")
					reckon.That(func() {
						reckon.That(map[int]string{3: "a"}).Has.No.Property(3, "a")
					}).Will.PanicWith("Does have property")
					reckon.That(func() {
						reckon.That(map[string]int{"a": 3}).Has.No.Property("a")
					}).Will.PanicWith("Does have property")
					reckon.That(func() {
						reckon.That(map[string]int{"a": 3}).Has.No.Property("a", 3)
					}).Will.PanicWith("Does have property")
					reckon.That(func() {
						reckon.That(struct{ A int }{A: 3}).Has.No.Property("A")
					}).Will.PanicWith("Does have property")
					reckon.That(func() {
						reckon.That(struct{ A int }{A: 3}).Has.No.Property("A", 3)
					}).Will.PanicWith("Does have property")
				})
			})
			suite.Describe("Deep Property", func(suite *suiteshop.Suite) {
				suite.Test("is", func(log *suiteshop.Log) {
					obj := map[string]interface{}{
						"a": struct {
							B interface{}
						}{
							B: []interface{}{
								map[string]interface{}{
									"c": true,
								},
							},
						},
					}
					reckon.That(obj).Has.DeepProperty("a.B.0.c")
					reckon.That(obj).Has.DeepProperty("a.B.0.c", true)
				})
			})
		})
	}).Post(fn) {
		t.Fatal(strings.Join(list, "\n"))
	} else {
		fmt.Println(strings.Join(list, "\n"))
	}
}
