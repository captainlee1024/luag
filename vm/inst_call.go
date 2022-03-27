package vm

import "github.com/captainlee1024/luag/api"

// 把当前函数的子函数原型实例化为闭包，并放入指定寄存器中
func closure(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1
	vm.LoadProto(bx)
	vm.Replace(a)
}

// 调用指定寄存器中的函数
func call(i Instruction, vm api.LuaVM) {
	// a func 调用函数在寄存器中的位置
	// b args 参数个数，在寄存器起始位置为a+1
	// c returns
	a, b, c := i.ABC()
	a += 1
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

// 把在指定寄存器中的函数和参数推到栈顶
func _pushFuncAndArgs(a, b int, vm api.LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		// 当已经有一部分参数停留在栈顶时，只需要把函数和前半部分参数
		// 推入栈顶，然后旋转栈顶即可
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm api.LuaVM) {
	// 获取栈顶记录上个调用帧返回值的数量
	// 目标寄存器索引
	x := int(vm.ToInteger(-1))
	// 弹出第一个值（上个调用帧执行结束后，返回值不弹出时，栈顶会记录下返回值个数）
	vm.Pop(1)

	vm.CheckStack(x - a)
	//
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	// 旋转栈调整函数和入参顺序
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func _popResults(a, c int, vm api.LuaVM) {
	if c == 1 {

	} else if c > 1 {
		// 栈顶是调用帧执行完毕之后留下的返回值
		// 替换到对应的寄存器中
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// 此时，保留下返回值，push原本要移动到目标寄存器的索引
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _return(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	if b == 1 {
		// no return values
	} else if b > 1 {
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else { // B == 0
		// 已经有一部分返回值在栈顶，把另一部分推入栈顶
		_fixStack(a, vm)
	}
}

func vararg(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	if b != 1 {
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

func tailCall(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

func self(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
func tForCall(i Instruction, vm api.LuaVM) {
	a, _, c := i.ABC()
	a += 1

	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_popResults(a+3, c+1, vm)
}
