package state

import (
	"github.com/captainlee1024/luag/api"

	"github.com/captainlee1024/luag/binchunk"
	"github.com/captainlee1024/luag/vm"
)

func (ls *luaState) Load(chunk []byte, chunkName, mod string) int {
	// TODO:编译源码还是解释字节码
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	ls.stack.push(c)
	if len(proto.Upvalues) > 0 { // 设置_ENV
		env := ls.registry.get(api.LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

func (ls *luaState) Call(nArgs, nResults int) {
	// 按照索引找到要调用的值,判断是否是函数
	val := ls.stack.get(-(nArgs + 1))
	c, ok := val.(*closure)

	if !ok { // 如果被调的值不是函数，则尝试执行元方法
		if mf := getMetafield(val, "__call", ls); mf != nil {
			if c, ok = mf.(*closure); ok {
				ls.stack.push(val)
				ls.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			//fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
			ls.callLuaClosure(nArgs, nResults, c)
		} else {
			ls.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function")
	}
}

func (ls *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, ls)
	newStack.closure = c

	if nArgs > 0 {
		args := ls.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	ls.stack.pop() // pop function

	//args := ls.stack.popN(nArgs)
	//newStack.pushN(args, nArgs)

	ls.pushLuaStack(newStack)
	r := c.goFunc(ls)
	ls.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		ls.stack.check(len(results))
		ls.stack.pushN(results, nResults)
	}
}

func (ls *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, ls)
	newStack.closure = c

	// 把方法和参数值一次性从栈顶弹出
	// 然后调用新帧的pushN方法,按照固定参数数量传入参数
	funcAndArgs := ls.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	// 修改新帧的栈顶指针
	newStack.top = nRegs
	// 记录 vararg参数
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// 新的调用帧推入栈顶
	ls.pushLuaStack(newStack)
	// 运行当前调用帧
	ls.runClosure()
	// 弹出运行完毕的调用帧
	ls.popLuaStack()

	// 此时，调用帧返回值停留在栈上
	// 弹出所有返回值
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		ls.stack.check(len(results))
		ls.stack.pushN(results, nResults)
	}
}

func (ls *luaState) runClosure() {
	for {
		inst := vm.Instruction(ls.Fetch())
		inst.Execute(ls)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

// Calls a function in protected mode.
// http://www.lua.org/manual/5.3/manual.html#lua_pcall
func (ls *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := ls.stack
	status = api.LUA_ERRRUN

	// catch error
	defer func() {
		if err := recover(); err != nil {
			if msgh != 0 {
				panic(err)
			}
			for ls.stack != caller {
				ls.popLuaStack()
			}
			ls.stack.push(err)
		}
	}()

	ls.Call(nArgs, nResults)
	status = api.LUA_OK
	return
}
