package state

func (state *luaState) PC() int {
	return state.pc
}

func (state *luaState) AddPC(n int) {
	state.pc += n
}

func (state *luaState) Fetch() uint32 {
	// TODO: 获取一条指令和运行该指令的详细过程
	i := state.proto.Code[state.pc]
	state.pc++
	return i
}

func (state *luaState) GetConst(idx int) {
	c := state.proto.Constants[idx]
	state.stack.push(c)
}

// GetRK 将指定常量或栈值推入栈顶
func (state *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant 常量
		state.GetConst(rk & 0xFF)
	} else { // register 寄存器
		state.PushValue(rk + 1)
	}
}
