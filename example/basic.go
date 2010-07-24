package main

import "lua51"
import "fmt"

func test(L *lua51.State) int {
	fmt.Println("hello world! from go!");
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


	lua51.Call(L,0,0);
	lua51.Call(L,0,0);
	lua51.Call(L,0,0);

	lua51.Close(L);
}
