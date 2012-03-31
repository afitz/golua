package lua

/*
#cgo pkg-config: lua5.1

#include <lua.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

type GoFunction func(*State) int;
type Alloc func(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer;

//wrapper to keep cgo from complaining about incomplete ptr type
//export State
type State struct {
	s *C.lua_State;
	//funcs []GoFunction;
	registry []interface{};
	//freelist for funcs indices, to allow for freeing
	freeIndices []uint;
}

//export golua_callgofunction
func golua_callgofunction(L interface{}, fid uint) int {
	L1 := L.(*State);
	f := L1.registry[fid].(GoFunction);
	return f(L1);
}

//export golua_gchook
func golua_gchook(L interface{}, id uint) int {
	L1 := L.(*State);
	L1.unregister(id);
	return 0;
}

//export golua_callpanicfunction
func golua_callpanicfunction(L interface{}, id uint) int {
	L1 := L.(*State);
	f := L1.registry[id].(GoFunction);
	return f(L1);
}

//export golua_idtointerface
func golua_idtointerface(id uint) interface{} {
	return id;
}

//export golua_cfunctiontointerface
func golua_cfunctiontointerface(f *uintptr) interface{} {
	return f;
}

//export golua_callallocf
func golua_callallocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	return uintptr((*((*Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr),osize,nsize));
}
