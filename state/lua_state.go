package state

// type luaState struct {
// 	stack *luaStack
// 	// 实现虚拟机新增的字段
// 	proto *binchunk.Prototype
// 	pc    int
// }

type luaState struct {
	stack *luaStack
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

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
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
