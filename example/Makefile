
include $(GOROOT)/src/Make.$(GOARCH)

all: basic panic

basic: basic.8
	$(LD) -o basic basic.8

basic.8: basic.go
	$(GC) -o basic.8 basic.go

panic: panic.8
	$(LD) -o panic panic.8

panic.8: panic.go
	$(GC) -o panic.8 panic.go

clean:
	-rm -f basic panic *.8
