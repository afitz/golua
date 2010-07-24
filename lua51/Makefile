# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

#ifndef TARG
#TARG=github.com/afitz/golua/lua5.1
#endif

CGO_CFLAGS+=
CGO_LDFLAGS+=_cgo_export.o -L./ -lgolua

ifndef LUA51_LIBNAME
LUA51_LIBNAME=lua5.1
endif

ifndef LUA51_INCLUDE_DIR
CGO_CFLAGS+=`pkg-config --cflags $(LUA51_LIBNAME)`
else
CGO_CFLAGS+=-I$(LUA51_INCLUDE_DIR)
endif

ifndef LUA51_LIB_DIR
CGO_LDFLAGS+=`pkg-config --libs $(LUA51_LIBNAME)`
else
CGO_LDFLAGS+=-L$(LUA51_LIB_DIR) -l$(LUA51_LIBNAME)
endif

TARG=lua51

#CGO_CFLAGS+=-I/usr/include/lua5.1
#CGO_LDFLAGS+=_cgo_export.o -L/usr/lib/lua5.1 -L./ -llua5.1 -lgolua
CGO_DEPS=libgolua.a _cgo_export.o
CGOFILES=\
	lua.go \
	lauxlib.go
GOFILES=lua_defs.go

CLEANFILES+=lua_defs.go

LUA_HEADER_FILES=/usr/include/lua5.1/lua.h /usr/include/lua5.1/lauxlib.h /usr/include/lua5.1/lualib.h

include $(GOROOT)/src/Make.pkg

%: install %.go
	$(QUOTED_GOBIN)/$(GC) $*.go
	$(QUOTED_GOBIN)/$(LD) -o $@ $*.$O

libgolua.a: golua.o
	ar rcs libgolua.a golua.o

golua.o: golua.c
	gcc -I/usr/include/lua5.1 -c golua.c -o golua.o

lua_defs.go:
	echo "package lua51;" > lua_defs.go
	cat $(LUA_HEADER_FILES) | grep '#define LUA' | sed 's/#define/const/' | sed 's/[A-Z_][A-Z_]*/& =/'  >> lua_defs.go
	
