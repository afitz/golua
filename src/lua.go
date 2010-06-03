package Lua

const LUA_VERSION	 = "Lua 5.1"
const LUA_RELEASE	 = "Lua 5.1.4"
const LUA_VERSION_NUM	 = 501
const LUA_COPYRIGHT	 = "Copyright (C) 1994-2008 Lua.org, PUC-Rio"
const LUA_AUTHORS	 = "R. Ierusalimschy, L. H. de Figueiredo & W. Celes"


/* mark for precompiled code (`<esc>Lua') */
const LUA_SIGNATURE	 = "\033Lua"

/* option for multiple returns in `lua_pcall' and `lua_call' */
const LUA_MULTRET	 = (-1)


/*
** pseudo-indices
*/
const LUA_REGISTRYINDEX	 = (-10000)
const LUA_ENVIRONINDEX	 = (-10001)
const LUA_GLOBALSINDEX	 = (-10002)
func lua_upvalueindex(i int) int { return LUA_GLOBALSINDEX-i; }


/* thread status; 0 is OK */
const LUA_YIELD	 = 1
const LUA_ERRRUN	 = 2
const LUA_ERRSYNTAX	 = 3
const LUA_ERRMEM	 = 4
const LUA_ERRERR	 = 5

type lua_GoFunction func(*lua_State) int;

type lua_Reader func(*lua_State, interface{}, *uintptr) string;
type lua_Writer func(*lua_State, *uintptr, uintptr, interface{}) int;

//declared in lmem.go because it already imports unsafe...
//type lua_Alloc func(interface{}, uintptr, uintptr, uintptr) uintptr;

/*
** basic types
*/
const LUA_TNONE		 = (-1);

const LUA_TNIL		 = 0
const LUA_TBOOLEAN		 = 1
const LUA_TLIGHTUSERDATA	 = 2
const LUA_TNUMBER		 = 3
const LUA_TSTRING		 = 4
const LUA_TTABLE		 = 5
const LUA_TFUNCTION		 = 6
const LUA_TUSERDATA		 = 7
const LUA_TTHREAD		 = 8



/* minimum Lua stack available to a C function */
const LUA_MINSTACK	 = 20

type lua_Number float32;
type lua_Integer int32;

//function prototypes removed

/*
** garbage-collection function and options
*/

const LUA_GCSTOP		 = 0
const LUA_GCRESTART		 = 1
const LUA_GCCOLLECT		 = 2
const LUA_GCCOUNT		 = 3
const LUA_GCCOUNTB		 = 4
const LUA_GCSTEP		 = 5
const LUA_GCSETPAUSE		 = 6
const LUA_GCSETSTEPMUL	 = 7


/* 
** ===============================================================
** some useful macros
** ===============================================================
*/

func lua_pop(L *lua_State,n int)		{
	lua_settop(L, -(n)-1);
}

func lua_newtable(L *lua_State)		{
	lua_createtable(L, 0, 0);
}

func lua_register(L *lua_State,n string,f *lua_GoFunction) {
	lua_pushcfunction(L, (f));
	lua_setglobal(L, (n));
}

func lua_pushcfunction(L *lua_State,f *lua_GoFunction)	{
	lua_pushcclosure(L, (f), 0);
}

func lua_strlen(L *lua_State,i int) uint {
	return lua_objlen(L, i);
}

func lua_isfunction(L *lua_State,n int)	 bool {
	return (lua_type(L, (n)) == LUA_TFUNCTION);
}
func lua_istable(L *lua_State,n int) bool {
	return (lua_type(L, (n)) == LUA_TTABLE);
}
func lua_islightuserdata(L *lua_State,n int) bool {
	return (lua_type(L, (n)) == LUA_TLIGHTUSERDATA);
}
func lua_isnil(L *lua_State,n int) bool	{
	return (lua_type(L, (n)) == LUA_TNIL);
}
func lua_isboolean(L *lua_State,n int) bool	{
	return (lua_type(L, (n)) == LUA_TBOOLEAN);
}
func lua_isthread(L *lua_State,n int) bool {
	return (lua_type(L, (n)) == LUA_TTHREAD);
}
func lua_isnone(L *lua_State,n int) bool {
	return (lua_type(L, (n)) == LUA_TNONE);
}
func lua_isnoneornil(L *lua_State, n int) bool {
	return (lua_type(L, (n)) <= 0);
}

//this may not be necessary in go....
func lua_pushliteral(L *lua_State, s string) {
	lua_pushlstring(L,s,len(s));
}

//#define lua_pushliteral(L, s)	\
//	lua_pushlstring(L, "" s, (sizeof(s)/sizeof(char))-1)

func lua_setglobal(L *lua_State,s string)	{
	lua_setfield(L, LUA_GLOBALSINDEX, (s));
}
func lua_getglobal(L *lua_State,s string)	{
	lua_getfield(L, LUA_GLOBALSINDEX, (s));
}

func lua_tostring(L *lua_State,i int) string {
	return lua_tolstring(L, (i), nil);
}



/*
** compatibility macros and functions
*/

func lua_open() *lua_State	{
	return luaL_newstate();
}

func lua_getregistry(L *lua_State)	{
	lua_pushvalue(L, LUA_REGISTRYINDEX);
}

func lua_getgccount(L *lua_State) int {
	return lua_gc(L, LUA_GCCOUNT, 0);
}

type lua_Chunkreader lua_Reader;
type lua_Chunkwriter lua_Writer;

//#define lua_Chunkreader		lua_Reader
//#define lua_Chunkwriter		lua_Writer

/*
** Event codes
*/
const LUA_HOOKCALL	 = 0
const LUA_HOOKRET	 = 1
const LUA_HOOKLINE	 = 2
const LUA_HOOKCOUNT	 = 3
const LUA_HOOKTAILRET  = 4


/*
** Event masks
*/
const LUA_MASKCALL	 = (1 << LUA_HOOKCALL)
const LUA_MASKRET	 = (1 << LUA_HOOKRET)
const LUA_MASKLINE	 = (1 << LUA_HOOKLINE)
const LUA_MASKCOUNT	 = (1 << LUA_HOOKCOUNT)

type lua_Hook func(L *lua_State, ar *lua_Debug) interface{};

type lua_Debug struct {
	event int;
	name string;
	namewhat string;
	what string;
	source string;
	currentline int;
	nups int;
	linedefined int;
	lastlinedefined int;
	short_src [LUA_IDSIZE]int8;
	/* private part */
	i_ci int;
};








