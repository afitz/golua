package main

import "../lua"
import "fmt"
import "errors"
import "os"

func testDefault(L *lua.State) {
	err := L.DoString("print(\"Unknown variable\" .. x)")
	fmt.Printf("Error is: %v\n", err)
	if err == nil {
		fmt.Printf("Error shouldn't have been nil\n")
		os.Exit(1)
	}
}

func faultyfunc(L *lua.State) int {
	panic(errors.New("An error"))
}

func faultyfunc2(L *lua.State) int {
	L.PushString("Some error")
	L.Error()
	return 1
}

func testRegistered(L *lua.State) {
	L.Register("faultyfunc", faultyfunc)
	err := L.DoString("faultyfunc()")
	fmt.Printf("Error is %v\n", err)
	if err == nil {
		fmt.Printf("Error shouldn't have been nil\n")
		os.Exit(1)
	}
}

func testRegistered2(L *lua.State) {
	L.Register("faultyfunc2", faultyfunc2)
	err := L.DoString("faultyfunc2()")
	fmt.Printf("Error is %v\n", err)
	if err == nil {
		fmt.Printf("Error shouldn't have been nil\n")
		os.Exit(1)
	}
}

func main() {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	testDefault(L)
	testRegistered(L)
	testRegistered2(L)
}
