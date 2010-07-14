#include <lua.h>
#include <lauxlib.h>
#include "_cgo_export.h"

//metatables to register:
//	GoLua.GoInterface
//  GoLua.GoFunction
//

static const char GoStateRegistryKey = 'k'; //golua registry key


void clua_initstate(lua_State* L)
{
	/* create the GoLua.GoFunction metatable */
	
}

unsigned int* clua_checkgofunction(lua_State* L, int index)
{
	unsigned int* fid = (unsigned int*)luaL_checkudata(L,index,"GoLua.GoFunction");
	luaL_argcheck(L, fid != NULL, index, "'GoFunction' expected");
	return fid;
}

//wrapper for callgofunction
int callback_function(lua_State* L)
{
	unsigned int *fid = clua_checkgofunction(L,-1);
	
	return golua_callgofunction(L,*fid);
}

unsigned int clua_togofunction(lua_State* L, int index)
{
	return *(clua_checkgofunction(L,index));
}

void clua_pushgofunction(lua_State* L, unsigned int fid)
{
	unsigned int* fidptr = (unsigned int*)lua_newuserdata(L, sizeof(unsigned int));
	*fidptr = fid;
	luaL_getmetatable(L, "GoLua.GoFunction");
	lua_setmetatable(L, -2);
}


void clua_setgostate(lua_State* L, GoInterface gi)
{
	lua_pushlightuserdata(L,(void*)&GoStateRegistryKey);
	GoInterface* gip = (GoInterface*)lua_newuserdata(L,sizeof(GoInterface));
	//copy interface value to userdata
	gip->v = gi.v;
	gip->t = gi.t;
	//set into registry table
	lua_settable(L,LUA_REGISTRYINDEX);
	
}

GoInterface clua_getgostate(lua_State* L)
{
	//get gostate from registry entry
	lua_pushlightuserdata(L,(void*)&GoStateRegistryKey);
	lua_gettable(L, LUA_REGISTRYINDEX);
	GoInterface* gip = lua_touserdata(L,-1);
	return *gip;	
}

void clua_pushgointerface(lua_State* L, GoInterface gi)
{
	GoInterface* iptr = (GoInterface*)lua_newuserdata(L, sizeof(GoInterface));
	iptr->v = gi.v;
	iptr->t = gi.t;
	luaL_getmetatable(L, "GoLua.GoInterface");
	lua_setmetatable(L,-2);
}
