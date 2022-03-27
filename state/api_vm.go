package state

func (ls *luaState) PC() int {
	return ls.stack.pc
}

func (ls *luaState) AddPC(n int) {
	ls.stack.pc += n
}

func (ls *luaState) Fetch() uint32 {
	// TODO: 获取一条指令和运行该指令的详细过程
	i := ls.stack.closure.proto.Code[ls.stack.pc]
	ls.stack.pc++
	return i
}

func (ls *luaState) GetConst(idx int) {
	c := ls.stack.closure.proto.Constants[idx]
	ls.stack.push(c)
}

// GetRK 将指定常量或栈值推入栈顶
func (ls *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant 常量
		ls.GetConst(rk & 0xFF)
	} else { // register 寄存器
		ls.PushValue(rk + 1)
	}
}

func (ls *luaState) RegisterCount() int {
	return int(ls.stack.closure.proto.MaxStackSize)
}

func (ls *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(ls.stack.varargs)
	}

	ls.stack.check(n)
	ls.stack.pushN(ls.stack.varargs, n)
}

func (ls *luaState) LoadProto(idx int) {
	/*
		proto := ls.stack.closure.proto.Protos[idx]
		closure := newLuaClosure(proto)
		ls.stack.push(closure)
	*/
	stack := ls.stack
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
func (ls *luaState) CloseUpvalues(a int) {
	for i, openuv := range ls.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(ls.stack.openuvs, i)
		}
	}
}
