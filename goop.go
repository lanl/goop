// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

/*
This package provides support for dynamic object-oriented programming
constructs in Go, much like those that appear in JavaScript.
*/
package goop

import "reflect"

// A goop.Object is a lot like a JavaScript object.
type Object struct {
	symbol_table map[string]interface{} // Map from a method name to a method value
}

// Create a new object with no methods defined on it.
func New() *Object {
	return &Object{symbol_table: make(map[string]interface{})}
}

// Assign a value to a member name.
func (obj Object) Set(memberName string, value interface{}) {
	obj.symbol_table[memberName] = value
}

// Return the value associated with a member name.
func (obj Object) Get(memberName string) (value interface{}) {
	value = obj.symbol_table[memberName]
	return
}

// Invoke a method on an object and return a slice of return values.
func (obj Object) Call(methodName string, arguments ...interface{}) []interface{} {
	// Construct a function and its arguments.
	userFuncIface := obj.symbol_table[methodName]
	userFunc := reflect.ValueOf(userFuncIface)
	userFuncArgs := make([]reflect.Value, len(arguments))
	for i, argIface := range arguments {
		userFuncArgs[i] = reflect.ValueOf(argIface)
	}

	// Call the function and return the results.
	returnVals := userFunc.Call(userFuncArgs)
	returnIfaces := make([]interface{}, len(returnVals))
	for i, val := range returnVals {
		returnIfaces[i] = val.Interface()
	}
	return returnIfaces
}
