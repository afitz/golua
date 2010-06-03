package Lua

//TODO: remove these placeholders
type StkId int;
type Closure struct{};
type lu_mem int;
//type GCObject struct{}
//type CommonHeader struct{};
type TValue struct{};
type lu_int32 int32;
type Instruction int;
type lu_byte byte;
type Mbuffer struct{};
type UpVal struct{};
type Table struct{};
//type lua_Hook func();
type TString struct{};

//const LUA_MINSTACK = 4;
const MAX_SIZET = 25;
//const NUM_TAGS = 8;
const TM_N = 8;
const LUA_IDSIZE = 4;

////////////////////////////////////
// lstate.h
////////////////////////////////////

//TODO: remove - forward declaration - defined in ldo.c
type lua_longjmp struct{}

func gt(L *lua_State) (*TValue) {
	return (L.l_gt);
}

func registry(L *lua_State) (*TValue) {
	return &(global(L).l_registry);
}

const EXTRA_STACK = 5;
const BASIC_CI_SIZE = 8;
const BASIC_STACK_SIZE = (2*LUA_MINSTACK);

type stringtable struct {
	hash []*GCObject;
	nuse lu_int32; //number of elements
	size int;
};

//information about a call
type CallInfo struct {
	base StkId;	//base for this function
	fun StkId;	//function index in stack
	top StkId;	//top for this function
	savedpc *Instruction; //TODO going to need to modify
						  //to make it actually work with the vm... i think
	nresults int; //expected num results from this function
	tailcalls int; //number of tail calls lost under this entry
};

func curr_func(L *lua_State) (*Closure) {
	return clvalue(L.ci.fun);
}

func ci_func(ci *CallInfo) (*Closure) {
	return clvalue(ci.fun);
}

func f_isLua(ci *CallInfo) bool {
	return !(ci_func(ci).c.isC);
}

func isLua(ci *CallInfo) bool {
	return ttisfunction(ci.fun) && f_isLua(ci);
}

/*
** `global state', shared by all threads of this state
*/
type global_State struct {
	strt *stringtable;  /* hash table for strings */
	frealloc lua_Alloc;  /* function to reallocate memory */
	ud interface{};         /* auxiliary data to `frealloc' */
	currentwhite lu_byte;
	gcstate lu_byte;  /* state of garbage collector */
	sweepstrgc int;  /* position of sweep in `strt' */
	rootgc []GCObject;  /* list of all collectable objects */
	//TODO: revisit the type here
	sweepgc **GCObject;  /* position of sweep in `rootgc' */
	gray []GCObject;  /* list of gray objects */
	grayagain []GCObject;  /* list of objects to be traversed atomically */
	weak []GCObject;  /* list of weak tables (to be cleared) */
	//TODO: revisit type of tmudata
	tmudata *GCObject;  /* last element of list of userdata to be GC */
	buff Mbuffer;  /* temporary buffer for string concatentation */
	GCthreshold lu_mem;
	totalbytes lu_mem;  /* number of bytes currently allocated */
	estimate lu_mem;  /* an estimate of number of bytes actually in use */
	gcdept lu_mem;  /* how much GC is `behind schedule' */
	gcpause int;  /* size of pause between successive GCs */
	gcstepmul int;  /* GC `granularity' */
	panicFunc lua_GoFunction;  /* to be called in unprotected errors */
	l_registry TValue;
	mainthread *lua_State;
	uvhead *UpVal;  /* head of double-linked list of all open upvalues */
	mt [NUM_TAGS]Table;  /* metatables for basic types */
	tmname [TM_N]TString;  /* array with tag-method names */
};

/*
** `per thread' state
*/
type lua_State struct {
    *CommonHeader;
	status lu_byte;
	top StkId;  /* first free slot in the stack */
	base StkId;  /* base of current function */
	l_G *global_State;
	ci *CallInfo;  /* call info for current function */
	savedpc *Instruction;  /* `savedpc' of current function */ //TODO: modify to work with actual vm implementation
	stack_last StkId;  /* last free slot in the stack */
	stack StkId;  /* stack base */
	end_ci *CallInfo;  /* points after end of ci array*/
	base_ci *CallInfo;  /* array of CallInfo's */
	stacksize int;
	size_ci int;  /* size of array `base_ci' */
	nCcalls int16;  /* number of nested C calls */
	baseCcalls int16;  /* nested C calls when resuming coroutine */
	hookmask lu_byte;
	allowhook lu_byte;
	basehookcount int;
	hookcount int;
	hook *lua_Hook;
	l_gt *TValue;  /* table of globals */
	env *TValue;  /* temporary place for environments */
	openupval []GCObject;  /* list of open upvalues in this stack */
	gclist []GCObject;
	errorJmp *lua_longjmp;  /* current error recover point */
	errfunc uintptr;  /* current error handling function (stack index) */
};

//'union' of all collectable objects
//TODO: might be able to eliminate GCHeader/COmmonheader of 
//	 	GC types, because no longer actually a union
type GCObject struct {
	gch *GCheader;
	value interface{};
}

//TODO: macros
//TODO: luaE_(new/free)thread
