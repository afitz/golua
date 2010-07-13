
include $(GOROOT)/src/Make.$(GOARCH)

basic: basic.8
	$(LD) -o basic basic.8

basic.8: basic.go
	$(GC) -o basic.8 basic.go

