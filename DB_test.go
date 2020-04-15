package main

import (
	"reflect"
	"testing"
)

func TestCheckUserReg(t *testing.T) {
	new_cases := []int{3, 5, 7, 9, 10}
	test_cases := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	expected := []bool{true, false, true, false, true, false, true, false, true, true}
	var got []bool

	for _, one_case := range new_cases {
		UserFirstStart(one_case)
	}

	for _, one_case := range test_cases {
		got = append(got, CheckUserReg(one_case))
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
