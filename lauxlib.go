package golua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include "golua.h"
import "C"
import "unsafe"
//luaL_addchar
//luaL_addlstring
//luaL_addsize
//luaL_addstring
//luaL_addvalue

func ArgCheck(L *State, cond bool, narg int, extramsg string) {
	if cond {
		C.luaL_argerror(L.s, C.int(narg), C.CString(extramsg));
	}
}

func ArgError(L *State, narg int, extramsg string) int {
	return int(C.luaL_argerror(L.s,C.int(narg),C.CString(extramsg)));
}

//type luaL_Buffer

//luaL_buffinit

func CallMeta(L *State, obj int, e string) int {
	return int(C.luaL_callmeta(L.s,C.int(obj),C.CString(e)));
}

func CheckAny(L *State, narg int) {
	C.luaL_checkany(L.s,C.int(narg));
}

func CheckInteger(L *State, narg int) int {
	return int(C.luaL_checkinteger(L.s,C.int(narg)));
}

func CheckNumber(L *State, narg int) float64 {
	return float64(C.luaL_checknumber(L.s,C.int(narg)));
}

func CheckString(L *State, narg int) string {
	var length C.size_t;
	return C.GoString( C.luaL_checklstring(L.s,C.int(narg),&length) );
}

func CheckOption(L *State, narg int, def string, lst []string) int {
	//TODO: complication: lst conversion to const char* lst[] from string slice
	return 0;
}

func CheckType(L *State, narg int, t int) {
	C.luaL_checktype(L.s,C.int(narg),C.int(t));
}
func CheckUdata(L *State, narg int, tname string) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_checkudata(L.s,C.int(narg),C.CString(tname)));
}

//true if no errors, false otherwise
func (L *State) DoFile(filename string) bool {
	if L.LoadFile(filename) == 0 {
		return L.PCall(0,LUA_MULTRET,0) == 0;
	}
	return false;
}

//true if no errors, false otherwise
func (L *State) DoString(str string) bool {
	if L.LoadString(str) == 0 {
		return L.PCall(0,LUA_MULTRET,0) == 0;
	}
	return false;
}

//luaL_error becomes FmtError because of lua_error
func FmtError(L *State, fmt string, args...interface{}) int {
	//TODO: complication: pass varargs
	return 0;
}

//returns false if no such metatable or no such field
func GetMetaField(L *State, obj int, e string) bool {
	return C.luaL_getmetafield(L.s,C.int(obj),C.CString(e)) != 0;
}

//TODO: rename better... clashes with lua_getmetatable
func LGetMetaTable(L *State, tname string) {
	//C.luaL_getmetatable(L.s,C.CString(tname));
	C.lua_getfield(L.s,LUA_REGISTRYINDEX,C.CString(tname));
}

func GSub(L *State, s string, p string, r string) string {
	return C.GoString(C.luaL_gsub(L.s, C.CString(s), C.CString(p), C.CString(r)));
}

//TODO: luaL_loadbuffer


func (L *State) LoadFile(filename string) int {
	return int(C.luaL_loadfile(L.s,C.CString(filename)));
}

func (L *State) LoadString(s string) int {
	return int(C.luaL_loadstring(L.s,C.CString(s)));
}

//returns false if registry already contains key tname
func (L *State) NewMetaTable(tname string) bool {
	return C.luaL_newmetatable(L.s, C.CString(tname)) != 0;
}

func NewState() *State {
	ls := (C.luaL_newstate());
	L := newState(ls);
	return L;
}

func (L *State) OpenLibs() {
	C.luaL_openlibs(L.s);
}

func (L *State) OptInteger(narg int, d int) int {
	return int(C.luaL_optinteger(L.s,C.int(narg),C.lua_Integer(d)));
}

func (L *State) OptNumber(narg int, d float64) float64 {
	return float64(C.luaL_optnumber(L.s,C.int(narg),C.lua_Number(d)));
}

func (L *State) OptString(narg int, d string) string {
	var length C.size_t;
	return C.GoString(C.luaL_optlstring(L.s,C.int(narg),C.CString(d),&length));
}

//luaL_prepbuffer

func (L *State) Ref(t int) int {
	return int(C.luaL_ref(L.s,C.int(t)));
}

//TODO: register - definately doable

//TODO: rename better
func LTypename(L *State, index int) string {
	return C.GoString(C.lua_typename(L.s,C.lua_type(L.s,C.int(index))));
}

//TODO: decide if we actually want this renamed
func TypeError(L *State, narg int, tname string) int {
	return int(C.luaL_typerror(L.s,C.int(narg),C.CString(tname)));
}

func Unref(L *State, t int, ref int) {
	C.luaL_unref(L.s,C.int(t),C.int(ref));
}

func Where(L *State, lvl int) {
	C.luaL_where(L.s,C.int(lvl));
}
