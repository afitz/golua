package Lua

/* tags for values visible from Lua */
const LAST_TAG	= LUA_TTHREAD;

const NUM_TAGS	= (LAST_TAG+1);


/*
** Extra tags for non-values
*/
const LUA_TPROTO	= (LAST_TAG+1);
const LUA_TUPVAL	= (LAST_TAG+2);
const LUA_TDEADKEY	= (LAST_TAG+3);

type CommonHeader struct {
	next *GCObject;
	tt lu_byte;
	marked lu_byte;
}

type GCheader CommonHeader;

//'union' of all Lua values
//here we use interface{}
type Value interface{};


type lua_TValue struct {
	value Value;
	tt int;
}

/* Macros to test type */
func ttisnil(o *lua_TValue) bool {
	return (ttype(o) == LUA_TNIL);
}
func ttisnumber(o *lua_TValue) bool {
	return (ttype(o) == LUA_TNUMBER);
}
func ttisstring(o *lua_TValue) bool {
	return (ttype(o) == LUA_TSTRING);
}
func ttistable(o *lua_TValue) bool {
	return (ttype(o) == LUA_TTABLE);
}
func ttisfunction(o *lua_TValue) bool {
	return (ttype(o) == LUA_TFUNCTION);
}
func ttisboolean(o *lua_TValue) bool {
	return (ttype(o) == LUA_TBOOLEAN);
}
func ttisuserdata(o *lua_TValue) bool {
	return (ttype(o) == LUA_TUSERDATA);
}
func ttisthread(o *lua_TValue) bool {
	return (ttype(o) == LUA_TTHREAD);
}
func ttislightuserdata(o *lua_TValue) bool {
	return (ttype(o) == LUA_TLIGHTUSERDATA);
}

/* Macros to access values */
func ttype(o *lua_TValue) int {
	return o.tt;
}

func gcvalue(o *lua_TValue)	*GCObject {
	return check_exp(iscollectable(o), (o.value).(*GCObject));
}
func pvalue(o *lua_TValue) uintptr {
	return check_exp(ttislightuserdata(o), (o.value).(uintptr));
}
func nvalue(o *lua_TValue) lua_Number {
	return check_exp(ttisnumber(o), (o.value).(lua_Number));
}
func rawtsvalue(o *lua_TValue) *TString	{
	return check_exp(ttisstring(o), &(((o.value).(*GCObject)).ts));
}
func tsvalue(o *lua_TValue) *TString {
	//return (&rawtsvalue(o)->tsv);
	return rawtsvalue(o);
}
func rawuvalue(o) *Udata {
	return check_exp(ttisuserdata(o), &(o)->value.gc->u);
}
func uvalue(o)	{
	return (&rawuvalue(o)->uv);
}
func clvalue(o)	{
	return check_exp(ttisfunction(o), &(o)->value.gc->cl);
}
func hvalue(o)	{
	return check_exp(ttistable(o), &(o)->value.gc->h);
}
func bvalue(o)	{
	return check_exp(ttisboolean(o), (o)->value.b);
}
func thvalue(o)	{
	return check_exp(ttisthread(o), &(o)->value.gc->th);
}
