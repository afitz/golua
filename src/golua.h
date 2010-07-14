//#include "_cgo_export.h"

typedef struct { void *t; void *v; } GoInterface;

unsigned int* clua_checkgofunction(lua_State* L, int index);
//int callback_function(lua_State* L)
unsigned int clua_togofunction(lua_State* L, int index);
void clua_pushgofunction(lua_State* L, unsigned int fid);
void clua_setgostate(lua_State* L, GoInterface gostate);
//void clua_pushgointerface(lua_State* L, GoInterface gi);
GoInterface clua_getgostate(lua_State* L);
