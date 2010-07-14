package lua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include "golua.h"
import "C"

func NewState() *State {
	ls := (C.luaL_newstate());
	L := newState(ls);
	return L;
}

func OpenLibs(L *State) {
	C.luaL_openlibs(L.s);
}

func DoString(L *State, str string) int {
	ok := C.luaL_loadstring(L.s, C.CString(str));
	if ok == 0 {
		return PCall(L,0,LUA_MULTRET, 0);
	}
	return int(ok);
}
