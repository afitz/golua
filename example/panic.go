package main

import lua "lua51"
import "fmt"

func test(L *lua.State) int {
	fmt.Println("hello world! from go!");
	return 0;
}

func main() {

	var L *lua.State;

	L = lua.NewState();
	lua.OpenLibs(L);

	currentPanicf := lua.AtPanic(L,nil);

	newPanic := func(L1 *lua.State) int {
		fmt.Println("I AM PANICKING!!!");
		return currentPanicf(L1);
	}

	lua.AtPanic(L,newPanic);

	//force a panic
	lua.PushNil(L);
	lua.Call(L,0,0);

	lua.Close(L);
}
