package main

import lua "lua51"
import "fmt"

// Simple execute code test
func test_simplestring(L *lua.State) {
	L.DoString(`
		print("test_simplestring: Hello, world!")
	`);
	// TODO: L.DoFile
}
// ---

// Simple callback - no input parameters, no output parameters
func test_callback_self(L *lua.State) int {
	fmt.Printf("test_callback_self: Callback ok\n")
	return 0
}

func test_callback_test(L *lua.State) {
	L.Register("test_callback", test_callback_self)
	L.DoString(`
		test_callback()
	`)
}
// ---

// Callback with input parameters
func test_cbinput_self(L *lua.State) int {
	param1 := L.ToString(1)
	param2 := L.ToInteger(2)
	fmt.Printf("test_cbinput_self: Callback with params: %s %d\n", param1, param2)
	return 0
}

func test_cbinput_test(L *lua.State) {
	L.Register("test_cbinput", test_cbinput_self)

	// call our function from lua
	L.DoString(`
		test_cbinput("test_string", 1234)
	`)

	// call our function from go - new variant
	L.GetField(lua.LUA_GLOBALSINDEX, "test_cbinput");
	L.PushAny("raw new")
	L.PushAny(2345)
	L.Call(2,0)

	// call our function from go - old variant
	L.GetField(lua.LUA_GLOBALSINDEX, "test_cbinput");
	L.PushString("raw old")
	L.PushInteger(3456)
	L.Call(2,0)
}
// ---

// Callback with output params
func test_cboutput_self(L *lua.State) int {
	L.PushAny("some string")
	L.PushAny(12345)
	return 2
}

func test_cboutput_run(L *lua.State) {
	L.Register("test_cboutput", test_cboutput_self)

	L.DoString(`
		p_str, p_int = test_cboutput()
		print("test_cboutput_self: p_str = ", p_str, " p_int = ", p_int)
	`)
}
// ---

func main() {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs();

	test_simplestring(L)
	test_callback_test(L)
	test_cbinput_test(L)
	test_cboutput_run(L)
}

