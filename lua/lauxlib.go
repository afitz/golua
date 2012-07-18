package lua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include <stdlib.h>
//#include "golua.h"
import "C"
import "unsafe"
//luaL_addchar
//luaL_addlstring
//luaL_addsize
//luaL_addstring
//luaL_addvalue

type LuaError struct {
	message string
}

func (err *LuaError) Error() string {
	return err.message
}

func ArgCheck(L *State, cond bool, narg int, extramsg string) {
	if cond {
		Cextramsg	:= C.CString(extramsg)
		defer C.free(unsafe.Pointer(Cextramsg))
		C.luaL_argerror(L.s, C.int(narg), Cextramsg)
	}
}

func ArgError(L *State, narg int, extramsg string) int {
	Cextramsg	:= C.CString(extramsg)
	defer C.free(unsafe.Pointer(Cextramsg))
	return int(C.luaL_argerror(L.s,C.int(narg),Cextramsg));
}

//type luaL_Buffer

//luaL_buffinit

func CallMeta(L *State, obj int, e string) int {
	Ce	:= C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return int(C.luaL_callmeta(L.s,C.int(obj),Ce))
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
	Ctname	:= C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_checkudata(L.s,C.int(narg),Ctname))
}

//true if no errors, false otherwise
func (L *State) DoFile(filename string) bool {
	if L.LoadFile(filename) == 0 {
		return L.PCall(0,LUA_MULTRET,0) == 0;
	}
	return false;
}

//nil if no errors, an error otherwise
func (L *State) DoString(str string) error {
	if L.LoadString(str) == 0 {
		if L.PCall(0,LUA_MULTRET,0) == 0 {
			return nil
		}
	}
	return &LuaError{L.ToString(-1)}
}

// evaluates argument like DoString, panics if execution failed
func (L *State) MustDoString(str string) {
	if err := L.DoString(str); err != nil {
		panic(err)
	}
}

//luaL_error becomes FmtError because of lua_error
func FmtError(L *State, fmt string, args...interface{}) int {
	//TODO: complication: pass varargs
	return 0;
}

//returns false if no such metatable or no such field
func GetMetaField(L *State, obj int, e string) bool {
	Ce	:= C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return C.luaL_getmetafield(L.s,C.int(obj),Ce) != 0;
}

//TODO: rename better... clashes with lua_getmetatable
func LGetMetaTable(L *State, tname string) {
	Ctname	:= C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	C.lua_getfield(L.s,LUA_REGISTRYINDEX,Ctname);
}

func GSub(L *State, s string, p string, r string) string {
	Cs	:= C.CString(s)
	Cp	:= C.CString(p)
	Cr	:= C.CString(r)
	defer func(){
		 C.free(unsafe.Pointer(Cs))
		 C.free(unsafe.Pointer(Cp))
		 C.free(unsafe.Pointer(Cr))
	}()

	return C.GoString(C.luaL_gsub(L.s, Cs, Cp, Cr))
}

//TODO: luaL_loadbuffer


func (L *State) LoadFile(filename string) int {
	Cfilename	:= C.CString(filename)
	defer C.free(unsafe.Pointer(Cfilename))
	return int(C.luaL_loadfile(L.s,Cfilename));
}

func (L *State) LoadString(s string) int {
	Cs	:= C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	return int(C.luaL_loadstring(L.s,Cs))
}

//returns false if registry already contains key tname
func (L *State) NewMetaTable(tname string) bool {
	Ctname	:= C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return C.luaL_newmetatable(L.s, Ctname) != 0;
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
	Cd := C.CString(d)
	defer C.free(unsafe.Pointer(Cd))
	return C.GoString(C.luaL_optlstring(L.s,C.int(narg),Cd,&length));
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
	Ctname	:= C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return int(C.luaL_typerror(L.s,C.int(narg),Ctname))
}

func Unref(L *State, t int, ref int) {
	C.luaL_unref(L.s,C.int(t),C.int(ref));
}

func Where(L *State, lvl int) {
	C.luaL_where(L.s,C.int(lvl));
}
