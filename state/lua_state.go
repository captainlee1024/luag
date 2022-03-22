package state

import "github.com/captainlee1024/luag/binchunk"

type luaState struct {
	stack *luaStack
	// 实现虚拟机新增的字段
	proto *binchunk.Prototype
	pc    int
}

/*
func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}
*/

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
