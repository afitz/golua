//#include "_cgo_export.h"

typedef struct { void *t; void *v; } GoInterface;

//function to setup metatables, etc
void clua_initstate(lua_State* L);

unsigned int clua_togofunction(lua_State* L, int index);
void clua_pushgofunction(lua_State* L, unsigned int fid);
void clua_setgostate(lua_State* L, GoInterface gostate);
GoInterface* clua_getgostate(lua_State* L);
GoInterface clua_atpanic(lua_State* L, unsigned int panicf_id);
int clua_callluacfunc(lua_State* L, lua_CFunction f);

//TODO: get/set panicf
//TODO: get/set allocf
//TODO: userdata support
