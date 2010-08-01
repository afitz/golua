
include $(GOROOT)/src/Make.$(GOARCH)

all: basic panic userdata

basic: basic.8
	$(LD) -o basic basic.8

basic.8: basic.go
	$(GC) -o basic.8 basic.go

panic: panic.8
	$(LD) -o panic panic.8

panic.8: panic.go
	$(GC) -o panic.8 panic.go

userdata: userdata.8
	$(LD) -o userdata userdata.8

userdata.8: userdata.go
	$(GC) -o userdata.8 userdata.go

clean:
	-rm -f basic panic *.8
