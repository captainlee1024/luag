package state

import "github.com/captainlee1024/luag/api"

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
func (ls *luaState) PushNil() {
	ls.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
func (ls *luaState) PushBoolean(b bool) {
	ls.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
func (ls *luaState) PushInteger(n int64) {
	ls.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
func (ls *luaState) PushNumber(n float64) {
	ls.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
func (ls *luaState) PushString(s string) {
	ls.stack.push(s)
}

// PushGoFunction 把go 函数转换成go闭包push入栈
func (ls *luaState) PushGoFunction(f api.GoFunction) {
	ls.stack.push(newGoClosure(f, 0))
}

// 弹出栈顶n个值作为GoClosure的Upvalue
// 然后把GoClosure推入栈顶
func (ls *luaState) PushGoClosure(f api.GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := ls.stack.pop()
		closure.upvals[i-1] = &upvalue{&val}
	}
}

// 全局变量表，也是一个普通的表，存放在registry中
// 在使用表API去操作全局变量表时，需要先把表push到栈里
func (ls *luaState) PushGlobalTable() {
	// TODO:?
	global := ls.registry.get(api.LUA_RIDX_GLOBALS)
	ls.stack.push(global)
}
