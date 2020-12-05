
typedef struct { void *t; void *v; } GoInterface;

/* function to setup metatables, etc */
void clua_initstate(lua_State* L);

unsigned int clua_togofunction(lua_State* L, int index);
void clua_pushgofunction(lua_State* L, unsigned int fid);
void clua_pushlightinteger(lua_State* L, int n);
int clua_tolightinteger(lua_State* L, int index);
void clua_setgostate(lua_State* L, GoInterface gostate);
GoInterface* clua_getgostate(lua_State* L);
GoInterface clua_atpanic(lua_State* L, unsigned int panicf_id);
int clua_callluacfunc(lua_State* L, lua_CFunction f);
lua_State* clua_newstate(void* goallocf);
void clua_setallocf(lua_State* L, void* goallocf);

void clua_openbase(lua_State* L);
void clua_openio(lua_State* L);
void clua_openmath(lua_State* L);
void clua_openpackage(lua_State* L);
void clua_openstring(lua_State* L);
void clua_opentable(lua_State* L);
void clua_openos(lua_State* L);


