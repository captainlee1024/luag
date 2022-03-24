package state

import "github.com/captainlee1024/luag/api"

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newtable
func (state *luaState) NewTable() {
	state.CreateTable(0, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_createtable
func (state *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	state.stack.push(t)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
func (state *luaState) GetTable(idx int) api.LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
func (state *luaState) GetField(idx int, k string) api.LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
func (state *luaState) GetI(idx int, i int64) api.LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, i)
}

// push(t[k])
func (state *luaState) getTable(t, k luaValue) api.LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		state.stack.push(v)
		return typeOf(v)
	}

	panic("not a table!") // todo
}
