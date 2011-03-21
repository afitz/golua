package main

import "golua"
import "fmt"

func test(L *golua.State) int {
	fmt.Println("hello world! from go!");
	return 0;
}

func test2(L *golua.State) int {
	arg := golua.CheckInteger(L,-1);
	argfrombottom := golua.CheckInteger(L,1);
	fmt.Print("test2 arg: ");
	fmt.Println(arg);
	fmt.Print("from bottom: ");
	fmt.Println(argfrombottom);
	return 0;
}

func main() {
	var L *golua.State;

	L = golua.NewState();
	L.OpenLibs();

	L.GetField(golua.LUA_GLOBALSINDEX, "print");
	L.PushString("Hello World!");
	L.Call(1,0);

	L.PushGoFunction(test);
	L.PushGoFunction(test);
	L.PushGoFunction(test);
	L.PushGoFunction(test);

	L.PushGoFunction(test2);
	L.PushInteger(42);
	L.Call(1,0);


	L.Call(0,0);
	L.Call(0,0);
	L.Call(0,0);

	L.Close();
}
