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

// Test constructors.
func TestConstructors(t *testing.T) {
	happyClass := func(self Object, someVal int) {
		self.Set("val", someVal)
	}
	value := 524287
	happyObj := New(happyClass, value)
	if result := happyObj.Get("val"); result.(int) != value {
		t.Fatalf("Expected %d but saw %v", value, result)
	}
}

// Test single-parent inheritance.
func TestSingleParentInheritance(t *testing.T) {
	// Test 1: Ensure that Get() finds members in an object's parent.
	point2D_class := func(self Object, x, y int) {
		// Construct a 2-D point.
		self.Set("x", x)
		self.Set("y", y)
	}
	point3D_class := func(self Object, x, y, z int) {
		// Construct a 3-D point.
		super := New(point2D_class, x, y)
		self.SetSuper(super)
		self.Set("z", z)
	}
	point3D := New(point3D_class, 2, 4, 8)
	expectedTotal := 2 + 4 + 8
	total := point3D.Get("x").(int) + point3D.Get("y").(int) + point3D.Get("z").(int)
	if total != expectedTotal {
		t.Fatalf("Expected %d from Get() but saw %d", expectedTotal, total)
	}

	// Test 2: Ensure that Contents() finds members in our parent.
	total = 0
	for key, value := range point3D.Contents(false) {
		switch key {
		case "x", "y", "z":
			total += value.(int)
		default:
			t.Fatalf("Did not expect key \"%s\", value %v", key, value)
		}
	}
	if total != expectedTotal {
		t.Fatalf("Expected %d from Contents() but saw %d", expectedTotal, total)
	}
}

// Test dynamically changing an object's lineage.
func TestSuperChange(t *testing.T) {
	parentType1 := func(self Object) {
		self.Set("one", 11111)
	}
	parentObj1 := New(parentType1)
	parentType2 := func(self Object) {
		self.Set("two", 22222)
	}
	parentObj2 := New(parentType2)
	childType := func(self Object) {
		self.SetSuper(parentObj1)
		self.Set("me", 33333)
	}
	childObj := New(childType)
	if result := childObj.Get("me").(int); result != 33333 {
		t.Fatalf("Expected %d but saw %v", 33333, result)
	}
	if result := childObj.Get("one").(int); result != 11111 {
		t.Fatalf("Expected %d but saw %v", 11111, result)
	}
	childObj.SetSuper(parentObj2)
	if result := childObj.Get("me").(int); result != 33333 {
		t.Fatalf("Expected %d but saw %v", 33333, result)
	}
	if result := childObj.Get("one"); result != NotFound {
		t.Fatalf("Expected %#v but saw %v", NotFound, result)
	}
	if result := childObj.Get("two").(int); result != 22222 {
		t.Fatalf("Expected %d but saw %v", 22222, result)
	}
	if result := childObj.GetSuper(); len(result) != 1 || !result[0].IsEquiv(parentObj2) {
		t.Fatalf("Expected equivalence between %#v and the first element of %#v", parentObj2, result)
	}
}

func TestDispatch(t *testing.T) {
	// Create an object with an "add" method that does different
	// things based on its arguments.
	adderObj := New()
	adderObj.Set("add", CombineFunctions(
		func(self Object, x, y int) int { return 10*x + y },
		func(self Object, x, y float64) float64 { return 100.0*x + y },
		func(self Object, a int) int { return -a }))

	// Test out the "add" method.
	if result := adderObj.Call("add", 77); result[0].(int) != -77 {
		t.Fatalf("Expected 77 but received %#v", result)
	}
	if result := adderObj.Call("add", 2, 3); result[0].(int) != 23 {
		t.Fatalf("Expected 23 but received %#v", result)
	}
	if result := adderObj.Call("add", 5.4, 3.2); result[0].(float64) != 543.2 {
		t.Fatalf("Expected 543.2 but received %#v", result)
	}
	if result := adderObj.Call("add", 5.4); result[0] != NotFound {
		t.Fatalf("Expected NotFound but received %#v", result)
	}
}
