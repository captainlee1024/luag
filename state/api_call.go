package state

import (
	"fmt"

	"github.com/captainlee1024/luag/binchunk"
	"github.com/captainlee1024/luag/vm"
)

func (state *luaState) Load(chunk []byte, chunkName, mod string) int {
	// TODO:编译源码还是解释字节码
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	state.stack.push(c)
	return 0
}

func (state *luaState) Call(nArgs, nResults int) {
	// 按照索引找到要调用的值,判断是否是函数
	val := state.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
		state.callLuaclosure(nArgs, nResults, c)
	} else {
		panic("not function")
	}
}

func (state *luaState) callLuaclosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVarage := (c.proto.IsVararg == 1)

	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	// 把方法和参数值一次性从栈顶弹出
	// 然后调用新帧的pushN方法,按照固定参数数量传入参数
	funcAndArgs := state.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	// 修改新帧的栈顶指针
	newStack.top = nRegs
	// 记录 varag参数
	if nArgs > nParams && isVarage {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// 新的调用帧推入栈顶
	state.pushLuaStack(newStack)
	// 运行当前调用帧
	state.runClosure()
	// 弹出运行完毕的调用帧
	state.popLuaStack()

	// 此时，调用帧返回值停留在栈上
	// 弹出所有返回值
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		state.stack.check(len(results))
		state.stack.pushN(results, nResults)
	}
}

func (state *luaState) runClosure() {
	for {
		inst := vm.Instruction(state.Fetch())
		inst.Execute(state)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}