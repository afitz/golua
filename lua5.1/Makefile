
include $(GOROOT)/src/Make.$(GOARCH)

GOTPATH=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)

OLDNAME=$(GOTPATH)/cgo_github.com_afitz_golua_lua5.1.so
NEWNAME=$(GOTPATH)/cgo_github.com_afitz_golua_lua51.so

all:
	make -f Make.lua all

clean:
	make -f Make.lua clean

install: $(NEWNAME)

$(NEWNAME): $(OLDNAME)
	cp $(GOTPATH)/cgo_github.com_afitz_golua_lua5.1.so $(GOTPATH)/cgo_github.com_afitz_golua_lua51.so

$(OLDNAME): *.c *.go
	make -f Make.lua install

uninstall:
	rm -f $(GOTPATH)/cgo_github.com_afitz_golua_lua51.so
	make -f Make.lua uninstall

