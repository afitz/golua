package lua51

//#include <lua.h>
//#include "golua.h"
import "C"

import "unsafe"
//TODO: remove
import "fmt"




//like lua_Writer, but as p will contain capacity, not needed as separate param
type Writer func(L *State, p []byte, ud interface{});
//like lua reader, but the return slice has the size, so does 
// we do not need it to be an out param
type Reader func(L *State, data interface{}) []byte;

//wrapper to keep cgo from complaining about incomplete ptr type
//export State
type State struct {
	s *C.lua_State;
	//funcs []GoFunction;
	registry []interface{};
	//freelist for funcs indices, to allow for freeing
	freeIndices []uint;
}

func newState(L *C.lua_State) *State {
	var newstatei interface{}
	newstate := &State{L, make([]interface{},0,8), make([]uint,0,8)};
	newstatei = newstate;
	ns1 := unsafe.Pointer(&newstatei);
	ns2 := (*C.GoInterface)(ns1);
	C.clua_setgostate(L,*ns2); //hacky....
	C.clua_initstate(L)
	return newstate;
}

func (L *State) addFreeIndex(i uint) {
	freelen := len(L.freeIndices)
	//reallocate if necessary
	if freelen+1 > cap(L.freeIndices) {
		newSlice := make([]uint, freelen, cap(L.freeIndices)*2);
		copy(newSlice, L.freeIndices);
		L.freeIndices = newSlice;
	}
	//reslice
	L.freeIndices = L.freeIndices[0:freelen+1];
	L.freeIndices[freelen] = i;
}

func (L *State) getFreeIndex() (index uint, ok bool) {
	freelen := len(L.freeIndices);
	//if there exist entries in the freelist
	if freelen > 0 {
		i := L.freeIndices[freelen - 1]; //get index
		L.freeIndices = L.freeIndices[0:i]; //'pop' index from list
		return i, true;
	}
	return 0,false;
}

//returns the registered function id
func (L *State) register(f interface{}) uint {
	index,ok := L.getFreeIndex();
	//if not ok, then we need to add new index by extending the slice
	if !ok {
		index = uint(len(L.registry));
		//reallocate backing array if necessary
		if index+1 > uint(cap(L.registry)) {
			newSlice := make([]interface{},index,cap(L.registry)*2);
			copy(newSlice, L.registry);
			L.registry = newSlice;
		}
		//reslice
		L.registry = L.registry[0:index+1]
	}
	L.registry[index] = f;
	return index;
}

func (L *State) unregister(fid uint) {
	if (fid < uint(len(L.registry))) && (L.registry[fid] != nil) {
		L.registry[fid] = nil;
		L.addFreeIndex(fid);
	}
}

type GoFunction func(*State) int;

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
	fmt.Printf("GC id: %d\n",id);
	return 0;
}

//export callpanicfunction
func callpanicfunction(L interface{}, id uint) int {
	L1 := L.(*State);
	f := L1.registry[id].(GoFunction);
	return f(L1);
}

//export idtointerface
func idtointerface(id uint) interface{} {
	return id;
}

//export cfunctiontointerface
func cfunctiontointerface(f *uintptr) interface{} {
	return f;
}


func (L *State) PushGoFunction(f GoFunction) {
	fid := L.register(f);
	C.clua_pushgofunction(L.s,C.uint(fid));
}


func (L *State) PushLightInteger(n int) {
	C.clua_pushlightinteger(L.s, C.int(n))
}


//push pointer by value, as a value - we don't impact lifetime
func (L *State) PushLightUserdata(ud *interface{}) {
	//push
	C.lua_pushlightuserdata(L.s,unsafe.Pointer(ud));

}
/*
//TODO:
//push pointer as full userdata - mem is go owned, but we 
//make a guarantee that lifetime will outlast lua lifetime
func PushUserdata(L *State, ud interface{}) {

}
*/

//old style
func (L *State) NewUserdata(size uintptr) unsafe.Pointer {
	return unsafe.Pointer(C.lua_newuserdata(L.s, C.size_t(size)));
}


type Alloc func(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer;
//export callAllocf
func callAllocf(fp uintptr,	ptr uintptr,
			    osize uint,			nsize uint) uintptr {
	return uintptr((*((*Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr),osize,nsize));
}

func (L *State) AtPanic(panicf GoFunction) (oldpanicf GoFunction) {
	fid := uint(0);
	if panicf != nil {
		fid = L.register(panicf);
	}
	oldres := interface{}(C.clua_atpanic(L.s,C.uint(fid)));
	switch i := oldres.(type) {
	case C.uint:
		f := L.registry[uint(i)].(GoFunction);
		//free registry entry
		L.unregister(uint(i));
		return f;
	case C.lua_CFunction:
		return func(L1 *State) int {
			return int(C.clua_callluacfunc(L1.s,i));
		}
	}
	//generally we only get here if the panicf got set to something like nil
	//potentially dangerous because we may silently fail 
	return nil;
}


func (L *State) Call(nargs int, nresults int) {
	C.lua_call(L.s,C.int(nargs),C.int(nresults));
}

func (L *State) CheckStack(extra int) bool {
	return C.lua_checkstack(L.s, C.int(extra)) != 0;
}

func (L *State) Close() {
	C.lua_close(L.s);
}

func (L *State) Concat(n int) {
	C.lua_concat(L.s,C.int(n));
}

func (L *State) CreateTable(narr int, nrec int) {
	C.lua_createtable(L.s, C.int(narr), C.int(nrec));
}

//CPcall replacement
func (L *State) GoPCall(fun GoFunction, ud interface{}) int {
	//TODO: need to emulate by pushing a c closure as in pushgofunction
	return 0;
}

//TODO: data be a slice? 
func (L *State) Dump( writer Writer, data interface{}) int {
	//TODO:
	return 0;
}

func (L *State) Equal(index1, index2 int) bool {
	return C.lua_equal(L.s, C.int(index1), C.int(index2)) == 1
}

func (L *State) Error() int	{ return int(C.lua_error(L.s)) }

func (L *State) GC(what, data int) int	{ return int(C.lua_gc(L.s, C.int(what), C.int(data))) }

func (L *State) GetfEnv(index int)	{ C.lua_getfenv(L.s, C.int(index)) }

func (L *State) GetField(index int, k string) {
	C.lua_getfield(L.s, C.int(index), C.CString(k))
}

func (L *State) GetGlobal(name string)	{ L.GetField(LUA_GLOBALSINDEX, name) }

func (L *State) GetMetaTable(index int) bool {
	return C.lua_getmetatable(L.s, C.int(index)) != 0
}

func (L *State) GetTable(index int)	{ C.lua_gettable(L.s, C.int(index)) }

func (L *State) GetTop() int	{ return int(C.lua_gettop(L.s)) }

func (L *State) Insert(index int)	{ C.lua_insert(L.s, C.int(index)) }

func (L *State) IsBoolean(index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == LUA_TBOOLEAN
}

func (L *State) IsGoFunction(index int) bool {
	//TODO:go function is now a userdatum, not a c function, so this will not work 
	return C.lua_iscfunction(L.s, C.int(index)) == 1
}

//TODO: add iscfunction

func (L *State) IsFunction(index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == LUA_TFUNCTION
}

func (L *State) IsLightUserdata(index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == LUA_TLIGHTUSERDATA
}

func (L *State) IsNil(index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == LUA_TNIL }

func (L *State) IsNone(index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == LUA_TNONE }

func (L *State) IsNoneOrNil(index int) bool	{ return int(C.lua_type(L.s, C.int(index))) <= 0 }

func (L *State) IsNumber(index int) bool	{ return C.lua_isnumber(L.s, C.int(index)) == 1 }

func (L *State) IsString(index int) bool	{ return C.lua_isstring(L.s, C.int(index)) == 1 }

func (L *State) IsTable(index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == LUA_TTABLE }

func (L *State) IsThread(index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == LUA_TTHREAD
}

func (L *State) IsUserdata(index int) bool	{ return C.lua_isuserdata(L.s, C.int(index)) == 1 }

func (L *State) LessThan(index1, index2 int) bool {
	return C.lua_lessthan(L.s, C.int(index1), C.int(index2)) == 1
}

func (L *State) Load(reader Reader, data interface{}, chunkname string) int {
	//TODO:
	return 0;
}

//NOTE: lua_newstate becomes NewStateAlloc whereas
//		luaL_newstate becomes NewState
func NewStateAlloc(f Alloc) *State {
	//TODO: implement a newState function which will initialize a State
	//		call with result from C.lua_newstate for the s initializer
	//ls := lua_newstate(
	ls := C.clua_newstate(unsafe.Pointer(&f));
	//ls := clua_newstate(
	return newState(ls);
}

func (L *State) NewTable() {
	C.lua_createtable(L.s,0,0);
}


func (L *State) NewThread() *State {
	//TODO: call newState with result from C.lua_newthread and return it
	//TODO: should have same lists as parent
	//		but may complicate gc
	s := C.lua_newthread(L.s);
	return &State{s,nil,nil};
}

func (L *State) Next(index int) int {
	return int(C.lua_next(L.s,C.int(index)));
}

func (L *State) ObjLen(index int) uint {
	return uint(C.lua_objlen(L.s,C.int(index)));
}

func (L *State) PCall(nargs int, nresults int, errfunc int) int {
	return int(C.lua_pcall(L.s, C.int(nargs), C.int(nresults), C.int(errfunc)));
}

func (L *State) Pop(n int) {
	//C.lua_pop(L.s, C.int(n));
	C.lua_settop(L.s, C.int(-n-1));
}

func (L *State) PushBoolean(b bool) {
	var bint int;
	if b {
		bint = 1;
	} else {
		bint = 0;
	}
	C.lua_pushboolean(L.s, C.int(bint));
}

func (L *State) PushString(str string) {
	C.lua_pushstring(L.s,C.CString(str));
}

func (L *State) PushInteger(n int) {
	C.lua_pushinteger(L.s,C.lua_Integer(n));
}

func (L *State) PushNil() {
	C.lua_pushnil(L.s);
}

func (L *State) PushNumber(n float64) {
	C.lua_pushnumber(L.s, C.lua_Number(n));
}

func (L *State) PushThread() (isMain bool) {
	return C.lua_pushthread(L.s) != 0;
}

func (L *State) PushValue(index int) {
	C.lua_pushvalue(L.s, C.int(index));
}

func (L *State) RawEqual(index1 int, index2 int) bool {
	return C.lua_rawequal(L.s, C.int(index1), C.int(index2)) != 0;
}

func (L *State) RawGet(index int) {
	C.lua_rawget(L.s, C.int(index));
}

func (L *State) RawGeti(index int, n int) {
	C.lua_rawgeti(L.s, C.int(index), C.int(n));
}

func (L *State) RawSet(index int) {
	C.lua_rawset(L.s, C.int(index));
}

func (L *State) RawSeti(index int, n int) {
	C.lua_rawseti(L.s, C.int(index), C.int(n));
}

func (L *State) Register(name string, f GoFunction) {
	L.PushGoFunction(f);
	L.SetGlobal(name);
}

func (L *State) Remove(index int) {
	C.lua_remove(L.s, C.int(index));
}

func (L *State) Replace(index int) {
	C.lua_replace(L.s, C.int(index));
}

func (L *State) Resume(narg int) int {
	return int(C.lua_resume(L.s, C.int(narg)));
}

func (L *State) SetAllocf(f Alloc) {
	C.clua_setallocf(L.s,unsafe.Pointer(&f));
}

func (L *State) SetfEnv(index int) {
	C.lua_setfenv(L.s, C.int(index));
}

func (L *State) SetField(index int, k string) {
	C.lua_setfield(L.s, C.int(index), C.CString(k));
}

func (L *State) SetGlobal(name string) {
	C.lua_setfield(L.s, C.int(LUA_GLOBALSINDEX), C.CString(name))
}

func (L *State) SetMetaTable(index int) {
	C.lua_setmetatable(L.s, C.int(index));
}

func (L *State) SetTable(index int) {
	C.lua_settable(L.s, C.int(index));
}

func (L *State) SetTop(index int) {
	C.lua_settop(L.s, C.int(index));
}

func (L *State) Status() int {
	return int(C.lua_status(L.s));
}

func (L *State) ToBoolean(index int) bool {
	return C.lua_toboolean(L.s, C.int(index)) != 0;
}

func (L *State) ToGoFunction(index int) (f GoFunction) {
	fid := C.clua_togofunction(L.s,C.int(index))
	return L.registry[fid].(GoFunction);
}

func (L *State) ToString(index int) string {
	var size C.size_t;
	//C.GoString(C.lua_tolstring(L.s, C.int(index), &size));
	return C.GoString(C.lua_tolstring(L.s,C.int(index),&size));
}

func (L *State) ToInteger(index int) int {
	return int(C.lua_tointeger(L.s, C.int(index)));
}

func (L *State) ToNumber(index int) float64 {
	return float64(C.lua_tonumber(L.s, C.int(index)));
}

func (L *State) ToPointer(index int) uintptr {
	return uintptr(C.lua_topointer(L.s, C.int(index)));
}

func (L *State) ToThread(index int) *State {
	//TODO: find a way to link lua_State* to existing *State, return that
	return &State{}
}

func (L *State) ToUserdata(index int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_touserdata(L.s,C.int(index)));
}

func (L *State) ToLightInteger(index int) int {
	return int(C.clua_tolightinteger(L.s, C.int(index)))
}

func (L *State) Type(index int) int {
	return int(C.lua_type(L.s, C.int(index)));
}

func (L *State) Typename(tp int) string {
	return C.GoString(C.lua_typename(L.s, C.int(tp)));
}

func XMove(from *State, to *State, n int) {
	C.lua_xmove(from.s, to.s, C.int(n));
}

func (L *State) Yield(nresults int) int {
	return int(C.lua_yield(L.s, C.int(nresults)));
}

// Restricted library opens

func (L *State) OpenBase() {
        C.clua_openbase(L.s);
}

func (L *State) OpenIO() {
        C.clua_openio(L.s);
}

func (L *State) OpenMath() {
        C.clua_openmath(L.s);
}

func (L *State) OpenPackage() {
        C.clua_openpackage(L.s);
}

func (L *State) OpenString() {
        C.clua_openstring(L.s);
}

func (L *State) OpenTable() {
        C.clua_opentable(L.s);
}

func (L *State) OpenOS() {
        C.clua_openos(L.s);
}
