// Copyright (C) 2011, Los Alamos National Security, LLC.
// Use of this source code is governed by a BSD-style license.

// This file ensures that the goop package is behaving itself properly.

package goop

import "testing"

func TestSimpleValues(t *testing.T) {
	value := 123
	obj := New()
	obj.Set("x", value)
	x := obj.Get("x").(int)
	if x != value {
		t.Fatalf("Expected x=%d but saw x=%v", value, x)
	}
}
