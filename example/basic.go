package main

import "lua"
import "fmt"

func main() {

	a := lua.GLOBALSINDEX;
	fmt.Println(a);

	var L *lua.State;

	lua.PushNumber(L,1.0);

}
