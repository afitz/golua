package lua

import (
	"testing"
)

type TestStruct struct {
	IntField int
	StringField string
	FloatField float64
}

func TestGoStruct(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	ts := &TestStruct{10, "test", 2.3 }

	L.CheckStack(1)

	L.PushGoStruct(ts)
	L.SetGlobal("t")

	L.GetGlobal("t")
	if !L.IsGoStruct(-1) { t.Fatal("Not go struct") }

	tsr := L.ToGoStruct(-1).(*TestStruct)
	if tsr != ts { t.Fatal("Retrieved something different from what we inserted") }

	L.Pop(1)

	L.PushString("This is not a struct")
	if L.ToGoStruct(-1) != nil {
		t.Fatal("Non-GoStruct value attempted to convert into GoStruct should result in nil")
	}

	L.Pop(1)
}


// TODO:
// - function test
// - struct field getting and setting test
// - DoFile / DoString tests (including errors

