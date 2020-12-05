# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc	

CGO_OFILES+=golua.o

ifndef LUA51_LIBNAME
LUA51_LIBNAME=lua5.1
endif

ifndef LUA51_INCLUDE_DIR
CGO_CFLAGS+=`pkg-config --cflags $(LUA51_LIBNAME)`
LUA51_INCLUDE_DIR:=$(shell pkg-config --cflags-only-I $(LUA51_LIBNAME) | sed 's/-I//' | sed 's/[ 	]*$$//')
else
CGO_CFLAGS+=-I$(LUA51_INCLUDE_DIR)
endif

ifndef LUA51_LIB_DIR
CGO_LDFLAGS+=`pkg-config --libs $(LUA51_LIBNAME)`
else
CGO_LDFLAGS+=-L$(LUA51_LIB_DIR) -l$(LUA51_LIBNAME)
endif

TARG=lua51

CGOFILES=\
	lua.go \
	lauxlib.go \
	lua_defs.go

CLEANFILES+=lua_defs.go

LUA_HEADERS=lua.h lauxlib.h lualib.h
LUA_HEADER_FILES:=$(patsubst %,$(LUA51_INCLUDE_DIR)/%,$(LUA_HEADERS))
LUA_INCLUDE_DIRECTIVES:=$(patsubst %,//\#include <%>\n, $(LUA_HEADERS))


include $(GOROOT)/src/Make.pkg

%: install %.go
	$(QUOTED_GOBIN)/$(GC) $*.go
	$(QUOTED_GOBIN)/$(LD) -o $@ $*.$O

golua.o: golua.c
	gcc $(CGO_CFLAGS) $(_CGO_CFLAGS_$(GOARCH)) -fPIC $(CFLAGS) -c golua.c -o golua.o

lua_defs.go:
	echo "package lua51;" > lua_defs.go
	echo "$(LUA_INCLUDE_DIRECTIVES)" "import \"C\"" >> lua_defs.go
#	echo "import \"C\"" >> lua_defs.go
	cat $(LUA_HEADER_FILES) | grep '#define LUA' | sed 's/#define/const/' | sed 's/\([A-Z_][A-Z_]*\)[\t ].*/\1 = C.\1/'  >> lua_defs.go

