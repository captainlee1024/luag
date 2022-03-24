package state

func (state *luaState) PC() int {
	return state.stack.pc
}

func (state *luaState) AddPC(n int) {
	state.stack.pc += n
}

func (state *luaState) Fetch() uint32 {
	// TODO: 获取一条指令和运行该指令的详细过程
	i := state.stack.closure.proto.Code[state.stack.pc]
	state.stack.pc++
	return i
}

func (state *luaState) GetConst(idx int) {
	c := state.stack.closure.proto.Constants[idx]
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

func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}

	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

func (self *luaState) LoadProto(idx int) {
	proto := self.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	self.stack.push(closure)
}
