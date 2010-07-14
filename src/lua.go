package lua

//#include <lua.h>
//#include "golua.h"
import "C"

import "unsafe"

//DEFINES
const GLOBALSINDEX = (-10002)

const (
	TNONE			= (-1)
	TNIL			= 0
	TBOOLEAN		= 1
	TLIGHTUSERDATA	= 2
	TNUMBER			= 3
	TSTRING			= 4
	TTABLE			= 5
	TFUNCTION		= 6
	TUSERDATA		= 7
	TTHREAD			= 8
)

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
	//todo add freelist for funcs indices, to allow for freeing
	freeIndices []uint;
}

func newState(L *C.lua_State) *State {
	var newstatei interface{}
	newstate := &State{L, make([]interface{},8), make([]uint,8)};
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
		newSlice := make([]uint, cap(L.freeIndices)*2);
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
			newSlice := make([]interface{},cap(L.registry)*2);
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
func golua_callgofunction(L *C.lua_State, fid uint) int {
	L1 := interface{}((C.clua_getgostate(L))).(*State);
	f := L1.registry[fid].(GoFunction);
	return f(L1);
}

//export golua_callpanicfunction
func callpanicfunction(L *C.lua_State) int {
	return 0;
}


func PushGoFunction(L *State, f GoFunction) {
	fid := L.register(f);
	C.clua_pushgofunction(L.s,C.uint(fid));
}

//TODO:
//push pointer by value, as a value - we don't impact lifetime
func PushLightUserdata(L *State, ud *interface{}) {
	//push
	C.lua_pushlightuserdata(L.s,unsafe.Pointer(ud));

}

//TODO:
//push pointer as full userdata - mem is go owned, but we 
//make a guarantee that lifetime will outlast lua lifetime
func PushUserdata(L *State, ud interface{}) {

}

//TODO:
//old style
func NewUserdata(L* State, size uintptr) uintptr {
	return 0;
}


type Alloc func(ud interface{}, ptr interface{}, osize uint, nsize uint) uintptr;

func AtPanic(L *State, panicf GoFunction) (oldpanicf GoFunction) {
	return *(new(GoFunction));
	//TODO: call lua_atpanic with c-wrapped panicf (may require field in State)
	//		check return of lua_atpanic to see if it is go or c function,
	//		if c, wrap as go.  if go, unwrap c portion
}


func Call(L *State, nargs int, nresults int) {
	C.lua_call(L.s,C.int(nargs),C.int(nresults));
}

func CheckStack(L *State, extra int) bool {
	return C.lua_checkstack(L.s, C.int(extra)) != 0;
}

func Close(L *State) {
	C.lua_close(L.s);
}

func Concat(L *State, n int) {
	C.lua_concat(L.s,C.int(n));
}

//CPcall replacement
func GoPCall(L *State, fun GoFunction, ud interface{}) int {
	//TODO: need to emulate by pushing a c closure as in pushgofunction
	return 0;
}

//TODO: data be a slice? 
func Dump(L *State, writer Writer, data interface{}) int {
	//TODO:
	return 0;
}

func Equal(L *State, index1, index2 int) bool {
	return C.lua_equal(L.s, C.int(index1), C.int(index2)) == 1
}

func Error(L *State) int	{ return int(C.lua_error(L.s)) }

func GC(L *State, what, data int) int	{ return int(C.lua_gc(L.s, C.int(what), C.int(data))) }

func GetfEnv(L *State, index int)	{ C.lua_getfenv(L.s, C.int(index)) }

func GetField(L *State, index int, k string) {
	C.lua_getfield(L.s, C.int(index), C.CString(k))
}

func GetGlobal(L *State, name string)	{ GetField(L, GLOBALSINDEX, name) }

func GetMetaTable(L *State, index int) bool {
	return C.lua_getmetatable(L.s, C.int(index)) != 0
}

func GetTable(L *State, index int)	{ C.lua_gettable(L.s, C.int(index)) }

func GetTop(L *State) int	{ return int(C.lua_gettop(L.s)) }

func Insert(L *State, index int)	{ C.lua_insert(L.s, C.int(index)) }

func IsBoolean(L *State, index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == TBOOLEAN
}

func IsGoFunction(L *State, index int) bool {
	//TODO: add a check if c function to distinguish c function from go function
	return C.lua_iscfunction(L.s, C.int(index)) == 1
}

//TODO: add iscfunction

func IsFunction(L *State, index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == TFUNCTION
}

func IsLightUserdata(L *State, index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == TLIGHTUSERDATA
}

func IsNil(L *State, index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == TNIL }

func IsNone(L *State, index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == TNONE }

func IsNoneOrNil(L *State, index int) bool	{ return int(C.lua_type(L.s, C.int(index))) <= 0 }

func IsNumber(L *State, index int) bool	{ return C.lua_isnumber(L.s, C.int(index)) == 1 }

func IsString(L *State, index int) bool	{ return C.lua_isstring(L.s, C.int(index)) == 1 }

func IsTable(L *State, index int) bool	{ return int(C.lua_type(L.s, C.int(index))) == TTABLE }

func IsThread(L *State, index int) bool {
	return int(C.lua_type(L.s, C.int(index))) == TTHREAD
}

func IsUserdata(L *State, index int) bool	{ return C.lua_isuserdata(L.s, C.int(index)) == 1 }

func LessThan(L *State, index1, index2 int) bool {
	return C.lua_lessthan(L.s, C.int(index1), C.int(index2)) == 1
}

func Load(L *State, reader Reader, data interface{}, chunkname string) int {
	//TODO:
	return 0;
}

func NewState(f Alloc, ud interface{}) *State {
	//TODO: implement a newState function which will initialize a State
	//		call with result from C.lua_newstate for the s initializer
	//ls := lua_newstate(
	
	return &State{};
}

func NewTable(L *State) {
	C.lua_createtable(L.s,0,0);
}

func NewThread(L* State) *State {
	//TODO: call newState with result from C.lua_newthread and return it
	//TODO: 
	s := C.lua_newthread(L.s);
	return &State{s,nil,nil};
}

func Next(L *State, index int) int {
	return int(C.lua_next(L.s,C.int(index)));
}

func ObjLen(L *State, index int) uint {
	return uint(C.lua_objlen(L.s,C.int(index)));
}

func PCall(L *State, nargs int, nresults int, errfunc int) int {
	return int(C.lua_pcall(L.s, C.int(nargs), C.int(nresults), C.int(errfunc)));
}

func Pop(L *State, n int) {
	//C.lua_pop(L.s, C.int(n));
	C.lua_settop(L.s, C.int(-n-1));
}

func PushBoolean(L *State, b bool) {
	var bint int;
	if b {
		bint = 1;
	} else {
		bint = 0;
	}
	C.lua_pushboolean(L.s, C.int(bint));
}

func PushString(L *State, str string) {
	//TODO:
}

func PushInteger(L *State, n int) {
	C.lua_pushinteger(L.s,C.lua_Integer(n));
}

func PushNil(L *State) {
	C.lua_pushnil(L.s);
}

func PushNumber(L *State, n float64) {
	C.lua_pushnumber(L.s, C.lua_Number(n));
}

func PushThread(L *State) (isMain bool) {
	return C.lua_pushthread(L.s) != 0;
}

func PushValue(L *State, index int) {
	C.lua_pushvalue(L.s, C.int(index));
}

func RawEqual(L *State, index1 int, index2 int) bool {
	return C.lua_rawequal(L.s, C.int(index1), C.int(index2)) != 0;
}

func RawGet(L *State, index int) {
	C.lua_rawget(L.s, C.int(index));
}

func RawGeti(L *State, index int, n int) {
	C.lua_rawgeti(L.s, C.int(index), C.int(n));
}

func RawSet(L *State, index int) {
	C.lua_rawset(L.s, C.int(index));
}

func RawSeti(L *State, index int, n int) {
	C.lua_rawseti(L.s, C.int(index), C.int(n));
}

func Register(L *State, name string, f GoFunction) {
	PushGoFunction(L,f);
	SetGlobal(L,name);
}

func Remove(L *State, index int) {
	C.lua_remove(L.s, C.int(index));
}

func Replace(L *State, index int) {
	C.lua_replace(L.s, C.int(index));
}

func Resume(L *State, narg int) int {
	return int(C.lua_resume(L.s, C.int(narg)));
}

func SetAllocf(L *State, f Alloc, ud interface{}) {
	//TODO:
}

func SetfEnv(L *State, index int) {
	C.lua_setfenv(L.s, C.int(index));
}

func SetField(L *State, index int, k string) {
	C.lua_setfield(L.s, C.int(index), C.CString(k));
}

func SetGlobal(L *State, name string) {
	C.lua_setfield(L.s, C.int(GLOBALSINDEX), C.CString(name))
}

func SetMetaTable(L *State, index int) {
	C.lua_setmetatable(L.s, C.int(index));
}

func SetTable(L *State, index int) {
	C.lua_settable(L.s, C.int(index));
}

func SetTop(L *State, index int) {
	C.lua_settop(L.s, C.int(index));
}

func Status(L *State) int {
	return int(C.lua_status(L.s));
}

func ToBoolean(L *State, index int) bool {
	return C.lua_toboolean(L.s, C.int(index)) != 0;
}

func ToGoFunction(L *State, index int) (f GoFunction) {
	fid := C.clua_togofunction(L.s,C.int(index))
	return L.registry[fid].(GoFunction);
}

func ToString(L *State, index int) string {
	var size C.size_t;
	return C.GoString(C.lua_tolstring(L.s, C.int(index), &size));
}

func ToInteger(L *State, index int) int {
	return int(C.lua_tointeger(L.s, C.int(index)));
}

func ToNumber(L *State, index int) float64 {
	return float64(C.lua_tonumber(L.s, C.int(index)));
}

func ToPointer(L *State, index int) uintptr {
	return uintptr(C.lua_topointer(L.s, C.int(index)));
}

func ToThread(L *State, index int) *State {
	//TODO: find a way to link lua_State* to existing *State, return that
	return &State{}
}

func ToUserdata(L *State, index int) interface{} {
	//TODO: needs userdata implementation first...
	return 0;
}

func Type(L *State, index int) int {
	return int(C.lua_type(L.s, C.int(index)));
}

func Typename(L *State, tp int) string {
	return C.GoString(C.lua_typename(L.s, C.int(tp)));
}

func XMove(from *State, to *State, n int) {
	C.lua_xmove(from.s, to.s, C.int(n));
}

func Yield(L *State, nresults int) int {
	return int(C.lua_yield(L.s, C.int(nresults)));
}



