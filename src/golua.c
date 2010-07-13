#include <lua.h>
#include <lauxlib.h>
#include "_cgo_export.h"

int callback_function(lua_State* L)
{
	int fid = lua_tointeger(L, lua_upvalueindex(1));
	
	//expecting that the global __GOLUA_STATE is a userdata filled with:
	// and GoState* cast to int*, so the pointer to the userdata is int**
	lua_getglobal(L, "___GOLUA_STATE");
	int** goState = (int**)(lua_topointer(L,-1));
	return golua_callgofunction(*goState,fid);
}

void clua_pushgofunction(lua_State* L, unsigned int fid)
{
	lua_pushinteger(L, fid);
	lua_pushcclosure(L, &callback_function, 1);
}

void clua_setgostate(lua_State* L, int* gostate)
{
	int** gostateptr = (int**)lua_newuserdata(L,sizeof(int*));
	*gostateptr = gostate;
	lua_registerglobal(L,-1,"___GOLUA_STATE");  //TODO: special stack id?
}

void clua_pushgointerface(lua_State* L, GoInterface gi)
{
	GoInterface* iptr = (GoInterface*)lua_newuserdata(L, sizeof(GoInterface));
	iptr->v = gi.v;
	iptr->t = gi.t;
	luaL_getmetatable(L, "GoLua.GoInterface");
	lua_setmetatable(L,-2);
}
