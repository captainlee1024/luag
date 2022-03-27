package state

import "github.com/captainlee1024/luag/api"

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (state *luaState) SetTable(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (state *luaState) SetField(idx int, k string) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (state *luaState) SetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v, false)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, true)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, self)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, self)
	} else {
		panic("table expected!") // todo
	}
}

// t[k]=v
func (state *luaState) setTable(t, k, v luaValue, raw bool) {
	//if tbl, ok := t.(*luaTable); ok {
	//	tbl.put(k, v)
	//	return
	//}
	//
	//panic("not a table!")

	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	if !raw {
		if mf := getMetafield(t, "__newindex", state); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				state.setTable(x, k, v, false)
				return
			case *closure:
				state.stack.push(mf)
				state.stack.push(t)
				state.stack.push(k)
				state.stack.push(v)
				state.Call(3, 0)
				return
			}
		}
	}

	panic("index error!")
}

func (state *luaState) SetGlobal(name string) {
	t := state.registry.get(api.LUA_RIDX_GLOBALS)
	v := state.stack.pop()
	state.setTable(t, name, v, false)
}

func (state *luaState) Register(name string, f api.GoFunction) {
	state.PushGoFunction(f)
	state.SetGlobal(name)
}
