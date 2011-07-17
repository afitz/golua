# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

LUA_INCLUDE_DIR=/usr/include/lua5.1/
LUA_LIB_DIR=/usr/lib/

CGO_CFLAGS+=-I$(LUA_INCLUDE_DIR)
CGO_DEPS=_cgo_export.o
CGO_LDFLAGS+=-L$(LUA_LIB_DIR) -lm -llua5.1
CGO_OFILES=\
		   golua.o \

TARG=golua

CGOFILES=\
	lua.go \
	lauxlib.go \
	lua_defs.go

CLEANFILES+=lua-5.1.4/src/*.o\
			example/*.8\
			example/basic\
			example/alloc\
			example/panic\
			example/userdata\

LUA_HEADERS=lua.h lauxlib.h lualib.h
LUA_HEADER_FILES:=$(patsubst %,$(LUA_INCLUDE_DIR)%,$(LUA_HEADERS))
LUA_INCLUDE_DIRECTIVES:=$(patsubst %,//\#include <%>\n, $(LUA_HEADERS))

include $(GOROOT)/src/Make.pkg

all: install examples

%: install %.go
	$(QUOTED_GOBIN)/$(GC) $*.go
	$(QUOTED_GOBIN)/$(LD) -o $@ $*.$O

golua.o: golua.c
	gcc $(CGO_CFLAGS) $(_CGO_CFLAGS_$(GOARCH)) -fPIC $(CFLAGS) -c golua.c -o golua.o

genluadefs:
	echo "package golua;" > lua_defs.go
	echo "$(LUA_INCLUDE_DIRECTIVES)" "import \"C\"" >> lua_defs.go
	echo "const (" >> lua_defs.go
	cat $(LUA_HEADER_FILES) | grep '#define LUA' | sed 's/#define/  /' | sed 's/\([A-Z_][A-Z_]*\)[[:space:]]*.*/\1 = C.\1/'  >> lua_defs.go
	echo ")" >> lua_defs.go

examples: install
	cd example && make
