// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

/*
This package provides support for dynamic object-oriented programming
constructs in Go, much like those that appear in JavaScript.
*/
package goop

import "reflect"

// An object is represented internally as a struct.
type internal struct {
	symbolTable map[string]interface{} // Map from a method name to a method value
}

// A goop.Error is used for producing return values that are
// differently typed from all user types.
type Error string

// A Get() call that fails to find the specified member returns goop.NotFound.
const NotFound = Error("Member not found")

// A goop.Object is a lot like a JavaScript object.
type Object struct {
	Implementation *internal // Internal representation not exposed to the user
}

// Allocate and return a new object.
func New() Object {
	obj := Object{}
	obj.Implementation = &internal{}
	obj.Implementation.symbolTable = make(map[string]interface{})
	return obj
}

// Assign a value to the name of an object member.
func (obj *Object) Set(memberName string, value interface{}) {
	obj.Implementation.symbolTable[memberName] = value
}

// Return the value associated with the name of an object member.
func (obj *Object) Get(memberName string) (value interface{}) {
	var ok bool
	if value, ok = obj.Implementation.symbolTable[memberName]; ok {
		return value
	}
	return NotFound
}

// Remove a member from an object.  This function always succeeds,
// even if the member did not previously exist.
func (obj *Object) Unset(memberName string) {
	obj.Implementation.symbolTable[memberName] = 0, false
}

// Return a map of all members of an object (useful for iteration).
// If the argument is true, also include method functions.
func (obj *Object) Contents(alsoMethods bool) map[string]interface{} {
	// Copy our internal structure to prevent the caller from
	// modifying it without our knowledge.
	symbolTable := obj.Implementation.symbolTable
	resultMap := make(map[string]interface{}, len(symbolTable))
	for key, val := range symbolTable {
		if alsoMethods || reflect.ValueOf(val).Kind() != reflect.Func {
			resultMap[key] = val
		}
	}
	return resultMap
}

// Invoke a method on an object and return the method's return values as a slice.
func (obj *Object) Call(methodName string, arguments ...interface{}) []interface{} {
	// Construct a function and its arguments.
	userFuncIface := obj.Implementation.symbolTable[methodName]
	userFunc := reflect.ValueOf(userFuncIface)
	userFuncArgs := make([]reflect.Value, len(arguments)+1)
	userFuncArgs[0] = reflect.ValueOf(*obj)
	for i, argIface := range arguments {
		userFuncArgs[i+1] = reflect.ValueOf(argIface)
	}

	// Call the function and return the results.
	returnVals := userFunc.Call(userFuncArgs)
	returnIfaces := make([]interface{}, len(returnVals))
	for i, val := range returnVals {
		returnIfaces[i] = val.Interface()
	}
	return returnIfaces
}
