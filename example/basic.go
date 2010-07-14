package main

import "lua"
import "fmt"

func test(L *lua.State) int {
	fmt.Println("hello world");
	return 0;
}

func main() {

	var L *lua.State;

	L = lua.NewState();
	lua.OpenLibs(L);

	lua.GetField(L, lua.LUA_GLOBALSINDEX, "print");
	lua.PushString(L, "Hello World!");
	lua.Call(L,1,0);

	lua.Close(L);
}
