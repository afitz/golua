package Lua

import "unsafe"
import "reflect"

/*****************
 * About 'vector' memory management:
 * - vectors are chunks of memory inserted as the backing for a slice 
 *	 object of appropriate type
 * - by convention the slice data is Len == Cap, because lua does 
 *   bookkeeping on its own about length already
 ******************/


////////////////////
// utility for this conversion
//////////////////////

func ptrElemSize(ptr interface{}) uintptr {
	v := reflect.NewValue(ptr);
	ve := reflect.Indirect(v);
	return ve.Type().Size();
}

type lua_Alloc func(ud interface{}, ptr unsafe.Pointer, oldsize uintptr, newsize uintptr) unsafe.Pointer;


func global(L *lua_State) *global_State {
	return new(global_State);
}

func luaD_throw(L *lua_State, err int) {

}

func luaG_runerror(L *lua_State, s string) {

}

///////////////////////////////////////////////////
// lmem.h conversion
///////////////////////////////////////////////////

const MEMERRMSG = "not enough memory";


func luaM_reallocv(L *lua_State,
		   block unsafe.Pointer,
		   oldNumElems uintptr,
		   newNumElems uintptr,
		   elemSize uintptr) unsafe.Pointer {
	if(newNumElems+1 <= uintptr(MAX_SIZET)/elemSize) {
		return luaM_realloc_(L,block,
				oldNumElems*elemSize,newNumElems*elemSize);
	}
	return luaM_toobig(L);
}

func luaM_freemem(L *lua_State,
		  block unsafe.Pointer, size uintptr) unsafe.Pointer {
	return luaM_realloc_(L,block,size,0);
}

func luaM_free(L* lua_State, ptr interface{} ) unsafe.Pointer {
	size := ptrElemSize(ptr);
	return luaM_realloc_(L,ptr.(unsafe.Pointer),size,0);
}

func luaM_freearray(L *lua_State, arr interface{}) unsafe.Pointer {
	//conv to sliceHeader
	arrType, arrPtr := unsafe.Reflect(arr);
	switch at := arrType.(type) {
	case *reflect.SliceType:
		elemSize := at.Elem().Size();
		var head reflect.SliceHeader
		head = unsafe.Unreflect(unsafe.Typeof(head), arrPtr).(reflect.SliceHeader);
		curSize := head.Cap;
		ptr := unsafe.Pointer(head.Data);
		return luaM_reallocv(L,ptr,uintptr(curSize),0,elemSize);
	default:
		//error
	}

	return nil;
}

func luaM_malloc(L *lua_State, size uintptr) unsafe.Pointer {
	return luaM_realloc_(L,nil,0,size);
}

func luaM_new(L* lua_State, ptr interface{}) interface{} {
	t := reflect.Typeof(ptr);
	switch pt := t.(type) {
	case *reflect.PtrType:
		// 'malloc' an object of type indicated by ptr
		elemT := pt.Elem();
		size := elemT.Size();
		newptr := luaM_malloc(L,size);
		//now cast it appropriately, so it can be .(type) cast
		castPtr := unsafe.Pointer(&newptr);
		return unsafe.Unreflect(unsafe.Typeof(ptr),castPtr);
	default:
		//error somehow
	}
	return nil;
}

//helper 
func ptr2slice(ptr unsafe.Pointer, st *reflect.SliceType, size int) interface{} {
	//now create the slice header, and unreflect to a slice
	sh := unsafe.Pointer(&reflect.SliceHeader{uintptr(ptr),size,size});
	return unsafe.Unreflect(st,sh);
}


//TODO: change to not use slicePtr but slice, maybe?
func luaM_newvector(L* lua_State, slicePtr interface{},
			numElems uintptr) interface{} {
	typ := reflect.Typeof(slicePtr);
	switch pt := typ.(type) {
	case *reflect.PtrType:
		var p interface{};
		p = pt.Elem();
		switch stt := p.(type) {
		case *reflect.SliceType:
			//allocate memory
			elemSize := uintptr(stt.Elem().Size());
			data := unsafe.Pointer(luaM_reallocv(L,nil,0,numElems,elemSize));
			return ptr2slice(data,stt,int(numElems));
		default:
			//error
		}
	default:
		//error
	}
	return nil;
}

func luaM_growvector(L *lua_State, slice interface{}, numElems uintptr, limit int, e string) interface{} {
	//conv to sliceHeader
	arrType, arrPtr := unsafe.Reflect(slice);
	switch at := arrType.(type) {
	case *reflect.SliceType:
		elemSize := at.Elem().Size();
		var head reflect.SliceHeader
		head = unsafe.Unreflect(unsafe.Typeof(head), arrPtr).(reflect.SliceHeader);
		curSize := int(head.Cap);
		ptr := unsafe.Pointer(head.Data);
		if(numElems+1 > uintptr(curSize)) {
			newPtr,newSize := luaM_growaux_(L,ptr,curSize,elemSize,limit,e);
			return ptr2slice(newPtr,at,newSize);
		}
	default:
		//error
	}
	return nil;
}



func luaM_reallocvector(L *lua_State, slice interface{}, newNumElem uintptr) interface{} {
	//conv to sliceHeader
	arrType, arrPtr := unsafe.Reflect(slice);
	switch at := arrType.(type) {
	case *reflect.SliceType:
		elemSize := at.Elem().Size();
		var head reflect.SliceHeader
		head = unsafe.Unreflect(unsafe.Typeof(head), arrPtr).(reflect.SliceHeader);
		curSize := head.Cap;
		ptr := unsafe.Pointer(head.Data);
		newPtr := luaM_reallocv(L,ptr,uintptr(curSize),newNumElem,elemSize);
		return ptr2slice(newPtr,at,int(newNumElem));
	default:
		//error
	}
	return nil;
}


////////////////////////////////////////////////////
// lmem.c conversion
////////////////////////////////////////////////////


const block_size = 100;

const MINSIZEARRAY = 4;

/* TODO: !!is going to require cgo! */
func l_alloc(ptr unsafe.Pointer, osize uintptr,nsize uintptr) unsafe.Pointer {
	//call C passthrough to malloc/realloc wrappers
	return nil;
}

func luaM_growaux_(L *lua_State,
		   block unsafe.Pointer,
		   size int, size_elems uintptr, limit int,
		   errormsg string) (unsafe.Pointer,int) {
	var newblock unsafe.Pointer;
	var newsize int;
	if(size >= limit/2) {
		if(size >= limit) {
			luaG_runerror(L,errormsg);
		}
		newsize = limit;
	} else {
		newsize = (size)*2;
		if(newsize < MINSIZEARRAY) {
			newsize = MINSIZEARRAY;
		}
	}
	newblock = luaM_reallocv(L, block, uintptr(size), uintptr(newsize), size_elems);
	return newblock,newsize;
}

func luaM_toobig(L *lua_State) unsafe.Pointer {
	luaG_runerror(L, "memory allocation error: block too big");
	return nil;
}

func luaM_realloc_(L *lua_State,
		   block unsafe.Pointer,
		   oldsize,size uintptr) unsafe.Pointer {
	g := global(L); //G(L)
	lua_assert((oldsize == 0) == (block == nil));
	block = g.frealloc(g.ud, block, oldsize, size);
	if(block == nil && size > 0) {
		luaD_throw(L, LUA_ERRMEM);
	}
	lua_assert((size == 0) == (block == nil));
	g.totalbytes = lu_mem(uintptr(g.totalbytes) - oldsize + size);
	return block;
}

