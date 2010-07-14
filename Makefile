
all: src/_obj/lua.a examples

src/_obj/lua.a:
	cd src && make

examples:
	cd example && make

clean:
	cd example && make clean
	cd src && make clean

install:
	cd src && make install
	
