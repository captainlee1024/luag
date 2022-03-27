package state

import "github.com/captainlee1024/luag/api"

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newtable
func (ls *luaState) NewTable() {
	ls.CreateTable(0, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_createtable
func (ls *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	ls.stack.push(t)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
func (ls *luaState) GetTable(idx int) api.LuaType {
	t := ls.stack.get(idx)
	k := ls.stack.pop()
	return ls.getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
func (ls *luaState) GetField(idx int, k string) api.LuaType {
	t := ls.stack.get(idx)
	return ls.getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
func (ls *luaState) GetI(idx int, i int64) api.LuaType {
	t := ls.stack.get(idx)
	return ls.getTable(t, i, false)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawget
func (ls *luaState) RawGet(idx int) api.LuaType {
	t := ls.stack.get(idx)
	k := ls.stack.pop()
	return ls.getTable(t, k, true)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgeti
func (ls *luaState) RawGetI(idx int, i int64) api.LuaType {
	t := ls.stack.get(idx)
	return ls.getTable(t, i, true)
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getmetatable
func (ls *luaState) GetMetatable(idx int) bool {
	val := ls.stack.get(idx)

	if mt := getMetatable(val, ls); mt != nil {
		ls.stack.push(mt)
		return true
	} else {
		return false
	}
}

// push(t[k])
func (ls *luaState) getTable(t, k luaValue, raw bool) api.LuaType {
	//if tbl, ok := t.(*luaTable); ok {
	//	v := tbl.get(k)
	//	ls.stack.push(v)
	//	return typeOf(v)
	//}
	//
	//
	//panic("not a table!")

	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetafield("__index") {
			ls.stack.push(v)
			return typeOf(v)
		}
	}

	if !raw {
		if mf := getMetafield(t, "__index", ls); mf != nil {
			switch x := mf.(type) {
			// 如果是表，则表的访问操作转发给该表，否则执行__index函数
			case *luaTable:
				return ls.getTable(x, k, false)
			case *closure:
				ls.stack.push(mf)
				ls.stack.push(t)
				ls.stack.push(k)
				ls.Call(2, 1)
				v := ls.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error!")
}

func (ls *luaState) GetGlobal(name string) api.LuaType {
	t := ls.registry.get(api.LUA_RIDX_GLOBALS)
	return ls.getTable(t, name, false)
}
