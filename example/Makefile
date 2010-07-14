
include $(GOROOT)/src/Make.$(GOARCH)

all: basic

basic: basic.8
	$(LD) -o basic basic.8

basic.8: basic.go
	$(GC) -o basic.8 basic.go


clean:
	-rm -f basic *.8
