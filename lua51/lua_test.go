package lua51

import "testing"
//import lua51 "lua51"

type luaFixture struct {
	L *State;
}

func TestGoFuncGC(t *testing.T) {
	L := NewState();
	OpenLibs(L);
	if len(L.freeIndices) != 0 {
		t.Logf("freeIndices != 0");
		t.Fail();
	}
	if len(L.registry) != 0 {
		t.Logf("registry != 0");
		t.Fail();
	}
	f := func (L *State) int {
		return 0;
	}
	PushGoFunction(L,f);
	if len(L.registry) != 1 {
		t.Logf("registry != 1");
		t.Fail();
	}
	if len(L.freeIndices) != 0 {
		t.Logf("freeIndices != 0");
		t.Fail();
	}
	Call(L,0,0);
	Close(L);
	if len(L.registry) != 1 {
		t.Logf("registry != 1");
		t.Fail();
	}
	if len(L.freeIndices) != 1 {
		t.Logf("freeIndices != 1");
		t.Fail();
	}
}
