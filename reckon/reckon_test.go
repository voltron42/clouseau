package reckon_test

import (
	"."
	"../suiteshop"

	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	fn := func(message string) {
		list = append(list, message)
	}
	suiteshop.Describe("Reckon", func(suite *suiteshop.Suite) {
		suite.Describe("Equals", func(suite *suiteshop.Suite) {
			suite.Test("is equal", func(log *suiteshop.Log) {
				reckon.That(5).Is.EqualTo(5)
			})
			suite.Test("not equal", func(log *suiteshop.Log) {
				reckon.That(4).Is.EqualTo(5)
			})
		})
		suite.Describe("Not Equal", func(suite *suiteshop.Suite) {
			suite.Test("is equal", func(log *suiteshop.Log) {
				reckon.That(5).Is.Not.EqualTo(5)
			})
			suite.Test("not equal", func(log *suiteshop.Log) {
				reckon.That(4).Is.Not.EqualTo(5)
			})
		})
		suite.Describe("Panic", func(suite *suiteshop.Suite) {
			suite.Test("will panic", func(log *suiteshop.Log) {
				reckon.That(func() {
					panic("panicking")
				}).Will.Panic()
			})
			suite.Test("will panic with wrong message", func(log *suiteshop.Log) {
				reckon.That(func() {
					panic("panicking")
				}).Will.PanicWith("should panic")
			})
		})
		suite.Describe("Not Panic", func(suite *suiteshop.Suite) {
			suite.Test("will panic", func(log *suiteshop.Log) {
				reckon.That(func() {
					panic("panicking")
				}).Will.Not.Panic()
			})
		})
	}).Post(fn)
	fmt.Println(strings.Join(list, "\n"))
}
