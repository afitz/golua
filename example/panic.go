package main

import lua "golua"
import "fmt"

func test(L *lua.State) int {
	fmt.Println("hello world! from go!");
	return 0;
}

func main() {

	var L *lua.State;

	L = lua.NewState();
	defer L.Close()
	L.OpenLibs();

	currentPanicf := L.AtPanic(nil);
	currentPanicf = L.AtPanic(currentPanicf);
	newPanic := func(L1 *lua.State) int {
		fmt.Println("I AM PANICKING!!!");
		return currentPanicf(L1);
	}

	L.AtPanic(newPanic);

	//force a panic
	L.PushNil();
	L.Call(0,0);
}
