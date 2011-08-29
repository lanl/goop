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
