// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file ensures that the goop package is behaving itself properly.

package goop

import "testing"

// Test setting and retrieving a scalar value.
func TestSimpleValues(t *testing.T) {
	value := 123
	var obj Object = New()
	obj.Set("x", value)
	x := obj.Get("x").(int)
	if x != value {
		t.Fatalf("Expected %d but saw %v", value, x)
	}
}

// Test multiply setting and retrieving a scalar value.
func TestMultipleAssign(t *testing.T) {
	value := 123
	var obj Object = New()
	obj.Set("x", value)
	value *= 2
	obj.Set("x", value)
	x := obj.Get("x").(int)
	if x != value {
		t.Fatalf("Expected %d but saw %v", value, x)
	}
}

// Test retrieving a nonexistent scalar value.
func TestNonexistent(t *testing.T) {
	var obj Object = New()
	x := obj.Get("bogus")
	if x != NotFound {
		t.Fatalf("Expectedly found member \"bogus\"")
	}
}

// Test creating and invoking a do-nothing method with no function
// arguments or return value.
func TestDoNothingFunction(t *testing.T) {
	var obj Object = New()
	obj.Set("doNothing", func(self Object) {})
	obj.Call("doNothing")
}

// Test invoking a method that returns its argument doubled.
func TestDoubleFunction(t *testing.T) {
	var obj Object = New()
	obj.Set("doubleIt", func(self Object, x int) int { return x * 2 })
	value := 123
	result := obj.Call("doubleIt", value)[0].(int)
	if result != value*2 {
		t.Fatalf("Expected %d but saw %v", value*2, result)
	}
}

// Test invoking a method that modifies object state.
func TestModifyObj(t *testing.T) {
	var obj Object = New()
	value := 100
	obj.Set("x", value)
	obj.Set("doubleX", func(self Object) {
		self.Set("x", self.Get("x").(int)*2)
	})
	obj.Call("doubleX")
	result := obj.Get("x").(int)
	if result != value*2 {
		t.Fatalf("Expected %d but saw %v", value*2, result)
	}
}

// Test iterating over all members.
func TestIteration(t *testing.T) {
	// Add various datatypes to an object.
	var obj Object = New()
	v_int := uint(2520)
	v_fp := float64(867.5309)
	v_str := "Yow!"
	obj.Set("integer", v_int)
	obj.Set("fp", v_fp)
	obj.Set("string", v_str)
	obj.Set("function", func() uint { return v_int })

	// Define a generic test for the above.
	test_contents := func(key string, value interface{}) {
		switch key {
		case "integer":
			if value.(uint) != v_int {
				t.Fatalf("Expected %d but saw %v", v_int, value)
			}
		case "fp":
			if value.(float64) != v_fp {
				t.Fatalf("Expected %.4f but saw %v", v_fp, value)
			}
		case "string":
			if value.(string) != v_str {
				t.Fatalf("Expected \"%s\" but saw %v", v_str, value)
			}
		default:
			t.Fatalf("Did not expect key \"%s\", value %v", key, value)
		}
	}

	// Test Contents(false).
	for key, value := range obj.Contents(false) {
		test_contents(key, value)
	}

	// Test Contents(true).
	for key, value := range obj.Contents(true) {
		if key == "function" {
			if funcResult := value.(func() uint)(); funcResult != v_int {
				t.Fatalf("Expected function \"%s\" to return %d, not %v", key, v_int, funcResult)
			}
		} else {
			test_contents(key, value)
		}
	}
}
