
LUA51_DIR=lua51

all: $(LUA51_DIR)/_obj/lua51.a examples

$(LUA51_DIR)/_obj/lua51.a:
	cd $(LUA51_DIR) && make 

examples: install
	cd example && make

clean:
	cd example && make clean
	cd $(LUA51_DIR) && make clean

install:
	cd $(LUA51_DIR) && make install
	
