// This file ensures that the goop package is behaving itself properly.

package goop_test

import (
	"fmt"
	"github.com/lanl/goop"
	"testing"
)

// Test setting and retrieving a scalar value.
func TestSimpleValues(t *testing.T) {
	value := 123
	obj := goop.New()
	obj.Set("x", value)
	x := obj.Get("x").(int)
	if x != value {
		t.Fatalf("Expected %d but saw %v", value, x)
	}
}

// Test multiply setting and retrieving a scalar value.
func TestMultipleAssign(t *testing.T) {
	value := 123
	obj := goop.New()
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
	obj := goop.New()
	x := obj.Get("bogus")
	if x != goop.ErrNotFound {
		t.Fatalf("Expectedly found member \"bogus\"")
	}
}

// Test creating and invoking a do-nothing method with no function
// arguments or return value.
func TestDoNothingFunction(t *testing.T) {
	obj := goop.New()
	obj.Set("doNothing", func(self goop.Object) {})
	obj.Call("doNothing")
}

// Test invoking a method that returns its argument doubled.
func TestDoubleFunction(t *testing.T) {
	obj := goop.New()
	obj.Set("doubleIt", func(self goop.Object, x int) int { return x * 2 })
	value := 123
	result := obj.Call("doubleIt", value)[0].(int)
	if result != value*2 {
		t.Fatalf("Expected %d but saw %v", value*2, result)
	}
}

// Test invoking a method that modifies object state.
func TestModifyObj(t *testing.T) {
	obj := goop.New()
	value := 100
	obj.Set("x", value)
	obj.Set("doubleX", func(self goop.Object) {
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
	obj := goop.New()
	vInt := uint(2520)
	vFp := float64(867.5309)
	vStr := "Yow!"
	obj.Set("integer", vInt)
	obj.Set("fp", vFp)
	obj.Set("string", vStr)
	obj.Set("function", func() uint { return vInt })

	// Define a generic test for the above.
	testContents := func(key string, value interface{}) {
		switch key {
		case "integer":
			if value.(uint) != vInt {
				t.Fatalf("Expected %d but saw %v", vInt, value)
			}
		case "fp":
			if value.(float64) != vFp {
				t.Fatalf("Expected %.4f but saw %v", vFp, value)
			}
		case "string":
			if value.(string) != vStr {
				t.Fatalf("Expected \"%s\" but saw %v", vStr, value)
			}
		default:
			t.Fatalf("Did not expect key \"%s\", value %v", key, value)
		}
	}

	// Test Contents(false).
	for key, value := range obj.Contents(false) {
		testContents(key, value)
	}

	// Test Contents(true).
	for key, value := range obj.Contents(true) {
		if key == "function" {
			if funcResult := value.(func() uint)(); funcResult != vInt {
				t.Fatalf("Expected function \"%s\" to return %d, not %v", key, vInt, funcResult)
			}
		} else {
			testContents(key, value)
		}
	}
}

// Test constructors.
func TestConstructors(t *testing.T) {
	happyClass := func(self goop.Object, someVal int) {
		self.Set("val", someVal)
	}
	value := 524287
	happyObj := goop.New(happyClass, value)
	if result := happyObj.Get("val"); result.(int) != value {
		t.Fatalf("Expected %d but saw %v", value, result)
	}
}

// Test single-parent inheritance.
func TestSingleParentInheritance(t *testing.T) {
	// Test 1: Ensure that Get() finds members in an object's parent.
	point2DClass := func(self goop.Object, x, y int) {
		// Construct a 2-D point.
		self.Set("x", x)
		self.Set("y", y)
	}
	point3DClass := func(self goop.Object, x, y, z int) {
		// Construct a 3-D point.
		super := goop.New(point2DClass, x, y)
		self.SetSuper(super)
		self.Set("z", z)
	}
	point3D := goop.New(point3DClass, 2, 4, 8)
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
	parentType1 := func(self goop.Object) {
		self.Set("one", 11111)
	}
	parentObj1 := goop.New(parentType1)
	parentType2 := func(self goop.Object) {
		self.Set("two", 22222)
	}
	parentObj2 := goop.New(parentType2)
	childType := func(self goop.Object) {
		self.SetSuper(parentObj1)
		self.Set("me", 33333)
	}
	childObj := goop.New(childType)
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
	if result := childObj.Get("one"); result != goop.ErrNotFound {
		t.Fatalf("Expected %#v but saw %v", goop.ErrNotFound, result)
	}
	if result := childObj.Get("two").(int); result != 22222 {
		t.Fatalf("Expected %d but saw %v", 22222, result)
	}
	if result := childObj.Super(); len(result) != 1 || !result[0].IsEquiv(parentObj2) {
		t.Fatalf("Expected equivalence between %#v and the first element of %#v", parentObj2, result)
	}
}

// Test the use of type-dependent dispatch (multiple methods with the
// same name but different types).
func TestDispatch(t *testing.T) {
	// Create an object with an "add" method that does different
	// things based on its arguments.
	adderObj := goop.New()
	adderObj.Set("add", goop.CombineFunctions(
		func(self goop.Object, x, y int) int { return 10*x + y },
		func(self goop.Object, x, y float64) float64 { return 100.0*x + y },
		func(self goop.Object, a int) int { return -a }))

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
	if result := adderObj.Call("add", 5.4); result[0] != goop.ErrNotFound {
		t.Fatalf("Expected ErrNotFound but received %#v", result)
	}
}

// The following is used by nativeFNV1.  We hope that making it
// exportable will prevent the compiler from optimizing it away.
var NativeHashVal uint64 = 14695981039346656037

// Apply the FNV-1 hash to a single 0xFF octet.
func nativeFNV1() {
	NativeHashVal *= 1099511628211
	NativeHashVal ^= 0xff
}

// Measure the speed of modifying a variable using native code.
func BenchmarkNativeFNV1(b *testing.B) {
	for i := b.N; i > 0; i-- {
		nativeFNV1()
	}
}

// Measure the speed of modifying a variable using native code.
func BenchmarkNativeFNV1Closure(b *testing.B) {
	b.StopTimer()
	var hashVal uint64 = 14695981039346656037
	fnv1 := func() {
		hashVal *= 1099511628211
		hashVal ^= 0xff
	}
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		fnv1()
	}
	b.StopTimer()
	if hashVal == 0 {
		// This case is pretty unlikely to occur.  We mainly
		// want to prevent hashVal from being optimized away.
		fmt.Printf("Cool!  We found a 64-bit zero hash of length %d.\n", b.N)
	}
}

// Measure the speed of modifying a variable using Goop's Get and Set methods.
func BenchmarkGoopFNV1(b *testing.B) {
	b.StopTimer()
	fnv1Obj := goop.New()
	fnv1Obj.Set("hashVal", uint64(14695981039346656037))
	fnv1 := func() {
		hashVal := fnv1Obj.Get("hashVal").(uint64)
		hashVal *= 1099511628211
		hashVal ^= 0xff
		fnv1Obj.Set("hashVal", hashVal)
	}
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		fnv1()
	}
}

// Measure the speed of modifying a variable using Goop's Get, Set,
// and Call methods.
func BenchmarkMoreGoopFNV1(b *testing.B) {
	b.StopTimer()
	fnv1Obj := goop.New()
	fnv1Obj.Set("hashVal", uint64(14695981039346656037))
	fnv1Obj.Set("fnv1", func(this goop.Object) {
		hashVal := this.Get("hashVal").(uint64)
		hashVal *= 1099511628211
		hashVal ^= 0xff
		this.Set("hashVal", hashVal)
	})
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		fnv1Obj.Call("fnv1")
	}
}

// Measure the speed of modifying a variable using Goop's Get, Set,
// Call, and CombineFunctions methods.
func BenchmarkEvenMoreGoopFNV1(b *testing.B) {
	b.StopTimer()
	fnv1Obj := goop.New()
	fnv1Obj.Set("hashVal", uint64(14695981039346656037))
	fnv1Obj.Set("fnv1", goop.CombineFunctions(
		func(this goop.Object) {
			hashVal := this.Get("hashVal").(uint64)
			hashVal *= 1099511628211
			hashVal ^= 0xff
			this.Set("hashVal", hashVal)
		}))
	b.StartTimer()
	for i := b.N; i > 0; i-- {
		fnv1Obj.Call("fnv1")
	}
}
