package main

import lua "lua51"
import "unsafe"
import "fmt"

type Userdata struct {
	a,b int
}

func main() {
	var L *lua.State;
	L = lua.NewState();
	lua.OpenLibs(L);

	rawptr := lua.NewUserdata(L,uintptr(unsafe.Sizeof(Userdata{})));
	var ptr *Userdata;
	ptr = (*Userdata)(rawptr);
	ptr.a = 2;
	ptr.b = 3;

	fmt.Println(ptr);

	rawptr2 := lua.ToUserdata(L,-1);
	ptr2 := (*Userdata)(rawptr2);

	fmt.Println(ptr2);
}
