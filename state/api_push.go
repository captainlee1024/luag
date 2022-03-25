package state

import "github.com/captainlee1024/luag/api"

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
func (state *luaState) PushNil() {
	state.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
func (state *luaState) PushBoolean(b bool) {
	state.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
func (state *luaState) PushInteger(n int64) {
	state.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
func (state *luaState) PushNumber(n float64) {
	state.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
func (state *luaState) PushString(s string) {
	state.stack.push(s)
}

// PushGoFunction 把go 函数转换成go闭包push入栈
func (state *luaState) PushGoFunction(f api.GoFunction) {
	state.stack.push(newGoClosure(f))
}

// 全局变量表，也是一个普通的表，存放在registry中
// 在使用表API去操作全局变量表时，需要先把表push到栈里
func (state *luaState) PushGlobalTable() {
	// TODO:?
	global := state.registry.get(api.LUA_RIDX_GLOBALS)
	state.stack.push(global)
}
