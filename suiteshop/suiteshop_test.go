package suiteshop_test

import (
	"."
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	list := []string{}
	suiteshop.Describe("Testing", func(suite *suiteshop.Suite) {
		suite.Test("1", func(log *suiteshop.Log) {
			log.Info("This is test 1.")
		})
		suite.Test("2", func(log *suiteshop.Log) {
			panic("This is test 2.")
		})
		suite.Test("3", func(log *suiteshop.Log) {
			log.Info("This is test 3.")
		})
		suite.Describe("deeper", func(suite *suiteshop.Suite) {
			suite.Test("1", func(log *suiteshop.Log) {
				log.Info("This is deeper test 1.")
			})
			suite.Test("2", func(log *suiteshop.Log) {
				panic("This is deeper test 2.")
			})
			suite.Test("3", func(log *suiteshop.Log) {
				log.Info("This is deeper test 3.")
			})
			suite.BeforeAll(func(log *suiteshop.Log) {
				log.Info("This happens before all the deeper tests.")
			})
			suite.AfterAll(func(log *suiteshop.Log) {
				log.Info("This happens after all the deeper tests.")
			})
			suite.Before(func(log *suiteshop.Log) {
				log.Info("This happens before each deeper test.")
			})
			suite.After(func(log *suiteshop.Log) {
				log.Info("This happens after each deeper test.")
			})
		})
		suite.Test("4", func(log *suiteshop.Log) {
			log.Info("This is test 4.")
		})
		suite.BeforeAll(func(log *suiteshop.Log) {
			log.Info("This happens before all the tests.")
		})
		suite.AfterAll(func(log *suiteshop.Log) {
			log.Info("This happens after all the tests.")
		})
		suite.Before(func(log *suiteshop.Log) {
			log.Info("This happens before each test.")
		})
		suite.After(func(log *suiteshop.Log) {
			log.Info("This happens after each test.")
		})
	}).Post(func(message string) {
		list = append(list, message)
	})
	fmt.Println(strings.Join(list, "\n"))
	if reflect.DeepEqual(list, []string{
		"This happens before all the tests.",
		"This happens before each test.",
		"This is test 1.",
		"This happens after each test.",
		"This happens before each test.",
		"This happens after each test.",
		"This happens before each test.",
		"This is test 3.",
		"This happens after each test.",
		"This happens after all the tests.",
		fmt.Sprintf("\x1b[42;37;1mTesting - 1"),
		fmt.Sprintf("\x1b[42;37;1mTesting - 2"),
		fmt.Sprintf("\x1b[41;37;1m\tTesting - 2 -- This is test 2."),
		fmt.Sprintf("\x1b[42;37;1mTesting - 3"),
	}) {
		t.Fatal("unequal")
	} else {
		fmt.Println("Safe!")
	}
}
