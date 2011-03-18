# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

CGO_CFLAGS+=-Ilua-5.1.4/src
CGO_DEPS=_cgo_export.o
CGO_LDFLAGS+=-lm
CGO_OFILES=\
		   golua.o \
		   lua-5.1.4/src/lapi.o \
		   lua-5.1.4/src/lauxlib.o \
		   lua-5.1.4/src/lbaselib.o \
		   lua-5.1.4/src/lcode.o \
		   lua-5.1.4/src/ldblib.o \
		   lua-5.1.4/src/ldebug.o \
		   lua-5.1.4/src/ldo.o \
		   lua-5.1.4/src/ldump.o \
		   lua-5.1.4/src/lfunc.o \
		   lua-5.1.4/src/lgc.o \
		   lua-5.1.4/src/linit.o \
		   lua-5.1.4/src/liolib.o \
		   lua-5.1.4/src/llex.o \
		   lua-5.1.4/src/lmathlib.o \
		   lua-5.1.4/src/lmem.o \
		   lua-5.1.4/src/loadlib.o \
		   lua-5.1.4/src/lobject.o \
		   lua-5.1.4/src/lopcodes.o \
		   lua-5.1.4/src/loslib.o \
		   lua-5.1.4/src/lparser.o \
		   lua-5.1.4/src/lstate.o \
		   lua-5.1.4/src/lstring.o \
		   lua-5.1.4/src/lstrlib.o \
		   lua-5.1.4/src/ltable.o \
		   lua-5.1.4/src/ltablib.o \
		   lua-5.1.4/src/ltm.o \
		   lua-5.1.4/src/lundump.o \
		   lua-5.1.4/src/lvm.o \
		   lua-5.1.4/src/lzio.o \
		   lua-5.1.4/src/print.o

TARG=golua

CGOFILES=\
	lua.go \
	lauxlib.go \
	lua_defs.go

CLEANFILES+=lua_defs.go

LUA_HEADERS=lua.h lauxlib.h lualib.h
LUA_HEADER_FILES:=$(patsubst %,lua-5.1.4/src/%,$(LUA_HEADERS))
LUA_INCLUDE_DIRECTIVES:=$(patsubst %,//\#include <%>\n, $(LUA_HEADERS))

include $(GOROOT)/src/Make.pkg

all: install examples

%: install %.go
	$(QUOTED_GOBIN)/$(GC) $*.go
	$(QUOTED_GOBIN)/$(LD) -o $@ $*.$O

golua.o: golua.c
	gcc $(CGO_CFLAGS) $(_CGO_CFLAGS_$(GOARCH)) -fPIC $(CFLAGS) -c golua.c -o golua.o

lua_defs.go:
	echo "package golua;" > lua_defs.go
	echo "$(LUA_INCLUDE_DIRECTIVES)" "import \"C\"" >> lua_defs.go
#	echo "import \"C\"" >> lua_defs.go
	cat $(LUA_HEADER_FILES) | grep '#define LUA' | sed 's/#define/const/' | sed 's/\([A-Z_][A-Z_]*\)[[:space:]]*.*/\1 = C.\1/'  >> lua_defs.go

examples: install
	cd example && make
