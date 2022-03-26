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

func (state *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(state.stack.varargs)
	}

	state.stack.check(n)
	state.stack.pushN(state.stack.varargs, n)
}

func (state *luaState) LoadProto(idx int) {
	/*
		proto := state.stack.closure.proto.Protos[idx]
		closure := newLuaClosure(proto)
		state.stack.push(closure)
	*/
	stack := state.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		// instack 在当前函数栈中，只需访问当前函数的局部变量
		if uvInfo.Instack == 1 {
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			// 处于open状态的，还在栈上，直接引用即可
			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else { // 处于闭合状态的需要保存在其他地方
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else { // 是更外层的变量，说明已经被当前函数的Upvalue捕获，从当前函数Upvalue获取即可
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

// 把open列表里的值copy一份，然后从openuv列表中删除
// TODO:?
func (state *luaState) CloseUpvalues(a int) {
	for i, openuv := range state.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(state.stack.openuvs, i)
		}
	}
}
