package mockband_test

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	if reflect.DeepEqual(5, 5) {
		fmt.Println("success")
	} else {
		t.Fatal("fail on identical")
	}
	if reflect.DeepEqual(5, 4) {
		t.Fatal("fail on number")
	} else {
		fmt.Println("success")
	}
	if reflect.DeepEqual(5, "a") {
		t.Fatal("fail on type")
	} else {
		fmt.Println("success")
	}
}
