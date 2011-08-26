// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file ensures that the goop package is behaving itself properly.

package goop

import "testing"

// Test setting and retrieving a scalar value.
func TestSimpleValues(t *testing.T) {
	value := 123
	var obj Object
	obj.Set("x", value)
	x := obj.Get("x").(int)
	if x != value {
		t.Fatalf("Expected %d but saw %v", value, x)
	}
}

// Test setting and invoking a do-nothing function with no arguments
// or return value.
func TestDoNothingFunction(t *testing.T) {
	var obj Object
	obj.Set("doNothing", func() {})
	obj.Call("doNothing")
}

// Test setting and invoking a function that returns its argument doubled.
func TestDoubleFunction(t *testing.T) {
	var obj Object
	obj.Set("doubleIt", func(x int) int {return x*2})
	value := 123
	result := obj.Call("doubleIt", value)[0].(int)
	if result != value*2 {
		t.Fatalf("Expected %d but saw %v", value*2, result)
	}
}
