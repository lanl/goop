/*
Package goop (Go Object-Oriented Programming) provides support for
dynamic object-oriented programming constructs in Go, much like those
that appear in various scripting languages.  The goal is to integrate
fast, native-Go objects and slower but more flexible Goop objects
within the same program.

FEATURES: For flexibility, Goop uses a prototype-based object model
(cf. http://en.wikipedia.org/wiki/Prototype-based_programming) rather
than a class-based object model.  Objects can be created either by
inheriting from existing objects or from scratch.  Data fields
(a.k.a. properties) and method functions can be added and removed at
will.  Multiple inheritance is supported.  An object's inheritance
hierarchy can be altered dynamically.  Methods can utilize
type-dependent dispatch (i.e., multiple methods with the same name but
different argument types).

As an example, let's create an object from scratch:

 pointObj := goop.New()

Now let's add a couple of data fields to pointObj:

 pointObj.Set("x", 0)
 pointObj.Set("y", 0)

Unlike native Go, Goop lets you define multiple method functions with
the same name, as long as the arguments differ in type and/or number:

 pointObj.Set("moveBy", goop.CombineFunctions(
         func(this goop.Object, xDelta, yDelta int) {
                 this.Set("x", this.Get("x") + xDelta)
                 this.Set("y", this.Get("y") + yDelta)
         },
         func(this goop.Object, delta int) {
                 this.Set("x", this.Get("x") + delta)
                 this.Set("y", this.Get("y") + delta)
         }))

Admittedly, having to use Get and Set all the time can be a bit
tedious.  Functions that are less trivial than the above will
typically call Get and Set only at the beginning and end of the
function and use local variables for most of the computation.

Use Call to call a method on an object:

 pointObj.Call("moveBy", 3, 5)
 pointObj.Call("moveBy", 12)

Call returns all of the method's return values as a single slice.  Use
type assertions to put the individual return values into their correct
format:

 pointObj.Set("distance", func(this goop.Object) float64 {
         x := float64(this.Get("x"))
         y := float64(this.Get("y"))
         return math.Sqrt(x*x + y*y)
 })

 d := pointObj.Call("distance")[0].(float64)

Again, sorry for the bloat, but that's what it takes to provide this
sort of dynamic behavior in Go.

The following more extended example shows how to define and
instantiate an LCMCalculator object, which is constructed from two
integers and provides methods that return the greatest common divisor
and least common multiple of those two numbers.  Each of those methods
memoizes its return value by redefining itself after its first
invocation to a function that returns a constant value.

 // This file showcases the Goop package by reimplementing the JavaScript LCM example from
 // http://en.wikipedia.org/wiki/Javascript#More_advanced_example.

 package main

 import "github.com/lanl/goop"
 import "fmt"
 import "sort"

 // Finds the lowest common multiple of two numbers
 func LCMCalculator(this goop.Object, x, y int) { // constructor function
         this.Set("a", x)
         this.Set("b", y)
         this.Set("gcd", func(this goop.Object) int { // method that calculates the greatest common divisor
                 abs := func(x int) int {
                         if x < 0 {
                                 x = -x
                         }
                         return x
                 }
                 a := abs(this.Get("a").(int))
                 b := abs(this.Get("b").(int))
                 if a < b {
                         // swap variables
                         a, b = b, a
                 }
                 for b != 0 {
                         t := b
                         b = a % b
                         a = t
                 }
                 // Only need to calculate GCD once, so "redefine" this
                 // method.  (Actually not redefinition - it's defined
                 // on the instance itself, so that this.gcd refers to
                 // this "redefinition".)
                 this.Set("gcd", func(this goop.Object) int { return a })
                 return a
         })
         this.Set("lcm", func(this goop.Object) int {
                 lcm := this.Get("a").(int) / this.Call("gcd")[0].(int) * this.Get("b").(int)
                 // Only need to calculate lcm once, so "redefine" this method.
                 this.Set("lcm", func(this goop.Object) int { return lcm })
                 return lcm
         })
         this.Set("toString", func(this goop.Object) string {
                 return fmt.Sprintf("LCMCalculator: a = %d, b = %d",
                         this.Get("a").(int), this.Get("b").(int))
         })
 }

 type lcmObjectVector []goop.Object

 func (lov lcmObjectVector) Less(i, j int) bool {
         a := lov[i].Call("lcm")[0].(int)
         b := lov[j].Call("lcm")[0].(int)
         return a < b
 }

 func (lov lcmObjectVector) Len() int {
         return len(lov)
 }

 func (lov lcmObjectVector) Swap(i, j int) {
         lov[i], lov[j] = lov[j], lov[i]
 }

 func main() {
         var lcmObjs lcmObjectVector
         for _, d := range [][]int{{25, 55}, {21, 56}, {22, 58}, {28, 56}} {
                 lcmObjs = append(lcmObjs, goop.New(LCMCalculator, d[0], d[1]))
         }
         sort.Sort(lcmObjs)
         for _, lcm := range lcmObjs {
                 fmt.Printf("%s, gcd = %d, lcm = %d\n",
                         lcm.Call("toString")[0], lcm.Call("gcd")[0], lcm.Call("lcm")[0])
         }
 }
*/
package goop

import "errors"
import "reflect"

// An object is represented internally as a struct.
type internal struct {
	symbolTable map[string]interface{} // Map from a member name to a member value
	prototypes  []Object               // List of other objects to search for members
}

// ErrNotFound is returned by a failed attempt to locate an object member.
var ErrNotFound = errors.New("Member not found")

// Object is a lot like a JavaScript object in that it uses prototype-based
// inheritance instead of a class hierarchy.
type Object struct {
	Implementation *internal // Internal representation not exposed to the user
}

// New allocates and return a new object.  It takes as arguments an
// optional constructor function with optional arguments.
func New(constructor ...interface{}) Object {
	// Allocate and initialize a new object.
	obj := Object{}
	obj.Implementation = &internal{}
	obj.Implementation.symbolTable = make(map[string]interface{})

	// If we weren't given a constructor, we have nothing left to
	// do.
	if len(constructor) == 0 {
		return obj
	}

	// Pass the new object and the given arguments to the
	// constructor.  Ignore the constructor's return value(s).
	constructorVal := reflect.ValueOf(constructor[0])
	argList := make([]reflect.Value, len(constructor))
	argList[0] = reflect.ValueOf(obj)
	for i, argIface := range constructor[1:] {
		argList[i+1] = reflect.ValueOf(argIface)
	}
	constructorVal.Call(argList)

	// Return the object we just constructed.
	return obj
}

// SetSuper specifies the object's parent object(s).  This is the
// mechanism by which both single and multiple inheritance are
// implemented.  For convenience, parents can be specified either
// individually or as a slice.
func (obj *Object) SetSuper(parentObjs ...interface{}) {
	// Empty the current set of prototypes.
	impl := obj.Implementation
	impl.prototypes = make([]Object, 0, len(parentObjs))

	// Append each prototype object in turn.
	for _, parentIface := range parentObjs {
		parentVal := reflect.ValueOf(parentIface)
		switch parentVal.Type().Kind() {
		case reflect.Array, reflect.Slice:
			// Append each object in turn to our prototype list.
			for i := 0; i < parentVal.Len(); i++ {
				impl.prototypes = append(impl.prototypes, parentVal.Index(i).Interface().(Object))
			}
		default:
			// Append the individual object to our prototype list.
			impl.prototypes = append(impl.prototypes, parentIface.(Object))
		}
	}
}

// Super returns the object's parent object(s) as a list.
func (obj *Object) Super() []Object {
	// Return a copy of impl.prototypes so if the caller mucks
	// with it, it won't mess up our object's internal
	// representation.
	protos := obj.Implementation.prototypes
	protoCopy := make([]Object, len(protos))
	copy(protoCopy, protos)
	return protoCopy
}

// IsEquiv returns whether another object is equivalent to the object
// in question.
func (obj *Object) IsEquiv(otherObj Object) bool {
	return obj.Implementation == otherObj.Implementation
}

// Set associates an arbitrary value with the name of an object member.
func (obj *Object) Set(memberName string, value interface{}) {
	obj.Implementation.symbolTable[memberName] = value
}

// Get returns the value associated with the name of an object member.
func (obj *Object) Get(memberName string) (value interface{}) {
	// Search our local members.
	var ok bool
	if value, ok = obj.Implementation.symbolTable[memberName]; ok {
		return value
	}

	// We didn't find the given member locally.  Try each of our
	// parents in turn.
	value = ErrNotFound
	for _, parent := range obj.Implementation.prototypes {
		parentValue := parent.Get(memberName)
		if parentValue != ErrNotFound {
			value = parentValue
			return
		}
	}
	return
}

// Unset removes a member from an object.  This function always
// succeeds, even if the member did not previously exist.
func (obj *Object) Unset(memberName string) {
	delete(obj.Implementation.symbolTable, memberName)
}

// Contents returns a map of all members of an object (useful for
// iteration).  If the argument is true, Contents also includes method
// functions.
func (obj *Object) Contents(alsoMethods bool) map[string]interface{} {
	// Copy our parents' data in reverse order so ancestor's
	// members are correctly overridden.
	impl := obj.Implementation
	resultMap := make(map[string]interface{}, len(impl.symbolTable))
	for i := len(impl.prototypes) - 1; i >= 0; i-- {
		parentObj := impl.prototypes[i]
		for key, val := range parentObj.Contents(alsoMethods) {
			resultMap[key] = val
		}
	}

	// Finally, copy our own object-specific data.
	for key, val := range impl.symbolTable {
		if alsoMethods || reflect.ValueOf(val).Kind() != reflect.Func {
			resultMap[key] = val
		}
	}
	return resultMap
}

// A typeDependentDispatch maps a textual type description to a
// function that accepts the associated types.
type typeDependentDispatch map[string]interface{}

// Given a function, functionSignature returns a string that describes
// its arguments.
func functionSignature(funcIface interface{}) string {
	funcType := reflect.ValueOf(funcIface).Type()
	numArgs := funcType.NumIn()
	argTypes := make([]byte, numArgs)
	for i := 0; i < numArgs; i++ {
		argTypes[i] = byte(funcType.In(i).Kind())
	}
	return string(argTypes)
}

// Given an array of arguments, argumentSignature returns a string
// that describes them.
func argumentSignature(argList []interface{}) string {
	numArgs := len(argList)
	argTypes := make([]byte, numArgs)
	for i := 0; i < numArgs; i++ {
		argTypes[i] = byte(reflect.TypeOf(argList[i]).Kind())
	}
	return string(argTypes)
}

// A MetaFunction encapsulates one or more functions, each with a
// unique argument-type signature.  When a MetaFunction is invoked, it
// accepts arbitrary inputs and returns arbitrary outputs (bundled
// into a slice).  On failure to find a matching signature, a
// singleton slice containing ErrNotFound is returned.
type MetaFunction func(varArgs ...interface{}) (funcResult []interface{})

// CombineFunctions combines multiple functions into a single
// MetaFunction for type-dependent dispatch.
func CombineFunctions(functions ...interface{}) MetaFunction {
	dispatchMap := make(typeDependentDispatch, len(functions))
	for _, funcIface := range functions {
		dispatchMap[functionSignature(funcIface)] = funcIface
	}
	dispatcher := func(varArgs ...interface{}) (funcResult []interface{}) {
		// Find the function in the dispatch map.
		funcIface, ok := dispatchMap[argumentSignature(varArgs)]
		if !ok {
			return []interface{}{ErrNotFound}
		}

		// Invoke the function.
		funcValue := reflect.ValueOf(funcIface)
		funcArgs := make([]reflect.Value, len(varArgs))
		for i, arg := range varArgs {
			funcArgs[i] = reflect.ValueOf(arg)
		}
		resultValues := funcValue.Call(funcArgs)

		// Convert the function's return values to a more
		// user-friendly type.
		funcResult = make([]interface{}, len(resultValues))
		for i, result := range resultValues {
			funcResult[i] = result.Interface()
		}
		return
	}
	return dispatcher
}

// Call invokes a method on an object and returns the method's return
// values as a slice.  Call returns a slice of the singleton ErrNotFound
// if the method could not be found.
func (obj *Object) Call(methodName string, arguments ...interface{}) []interface{} {
	// Construct a function and its arguments, using Get to
	// automatically search parent objects if necessary.
	userFuncIface := obj.Get(methodName)
	if userFuncIface == ErrNotFound {
		return []interface{}{ErrNotFound}
	}
	userFunc := reflect.ValueOf(userFuncIface)
	userFuncArgs := make([]reflect.Value, len(arguments)+1)
	userFuncArgs[0] = reflect.ValueOf(*obj)
	for i, argIface := range arguments {
		userFuncArgs[i+1] = reflect.ValueOf(argIface)
	}

	// Call the function.
	returnVals := userFunc.Call(userFuncArgs)
	returnIfaces := make([]interface{}, len(returnVals))
	for i, val := range returnVals {
		returnIfaces[i] = val.Interface()
	}

	// Return the results.  As a special case, we return a
	// MetaFunction's already-wrapped results without an
	// additional level of wrapping.
	if _, ok := userFuncIface.(MetaFunction); ok {
		returnIfaces = returnIfaces[0].([]interface{})
	}
	return returnIfaces
}
