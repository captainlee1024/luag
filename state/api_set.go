package state

import "github.com/captainlee1024/luag/api"

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (state *luaState) SetTable(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (state *luaState) SetField(idx int, k string) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (state *luaState) SetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v)
}

// t[k]=v
func (state *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}

	panic("not a table!")
}

func (state *luaState) SetGlobal(name string) {
	t := state.registry.get(api.LUA_RIDX_GLOBALS)
	v := state.stack.pop()
	state.setTable(t, name, v)
}

func (state *luaState) Register(name string, f api.GoFunction) {
	state.PushGoFunction(f)
	state.SetGlobal(name)
}
