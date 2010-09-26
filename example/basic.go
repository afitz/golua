package main

import "lua51"
import "fmt"

func test(L *lua51.State) int {
	fmt.Println("hello world! from go!");
	return 0;
}

func test2(L *lua51.State) int {
	arg := lua51.CheckInteger(L,-1);
	argfrombottom := lua51.CheckInteger(L,1);
	fmt.Print("test2 arg: ");
	fmt.Println(arg);
	fmt.Print("from bottom: ");
	fmt.Println(argfrombottom);
	return 0;
}

func main() {
	var L *lua51.State;

	L = lua51.NewState();
	lua51.OpenLibs(L);

	lua51.GetField(L, lua51.LUA_GLOBALSINDEX, "print");
	lua51.PushString(L, "Hello World!");
	lua51.Call(L,1,0);

	lua51.PushGoFunction(L, test);
	lua51.PushGoFunction(L, test);
	lua51.PushGoFunction(L, test);
	lua51.PushGoFunction(L, test);

	lua51.PushGoFunction(L, test2);
	lua51.PushInteger(L,42);
	lua51.Call(L,1,0);


	lua51.Call(L,0,0);
	lua51.Call(L,0,0);
	lua51.Call(L,0,0);

	lua51.Close(L);
}
