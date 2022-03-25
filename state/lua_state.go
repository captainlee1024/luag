package state

import "github.com/captainlee1024/luag/api"

// type luaState struct {
// 	stack *luaStack
// 	// 实现虚拟机新增的字段
// 	proto *binchunk.Prototype
// 	pc    int
// }

type luaState struct {
	registry *luaTable // 注册表
	stack    *luaStack
}

/*
func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}
*/

/*
func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
*/

/*
func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}
*/

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0)) // 全局环境

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(api.LUA_MINSTACK, ls))
	return ls
}

func (state *luaState) popLuaStack() {
	stack := state.stack
	state.stack = stack.prev
	stack.prev = nil
}

func (state *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = state.stack
	state.stack = stack
}
