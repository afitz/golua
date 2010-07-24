
all: lua5.1/_obj/lua5.1.a examples

lua5.1/_obj/lua5.1.a:
	cd lua5.1 && make

examples: install
	cd example && make

clean:
	cd example && make clean
	cd lua5.1 && make clean

install:
	cd lua5.1 && make install 
	
