package vm

import "github.com/captainlee1024/luag/api"

// R(A) := R(B)
func move(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a)
}

func jmp(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		// panic("todo: jmp!")
		vm.CloseUpvalues(a) // 局部变量声明周期结束，把处于open状态的引用转换为关闭状态
	}
}
