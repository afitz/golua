# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=lua
CGO_CFLAGS=-I/usr/include/lua5.1
CGO_LDFLAGS=_cgo_export.o -L/usr/lib/lua5.1 -L./ -llua5.1 -lgolua
CGO_DEPS=libgolua.a _cgo_export.o
CGOFILES=\
	lua.go

CLEANFILES+=

include $(GOROOT)/src/Make.pkg

%: install %.go
	$(QUOTED_GOBIN)/$(GC) $*.go
	$(QUOTED_GOBIN)/$(LD) -o $@ $*.$O

libgolua.a: golua.o
	ar rcs libgolua.a golua.o

golua.o: golua.c
	gcc -I/usr/include/lua5.1 -c golua.c -o golua.o
