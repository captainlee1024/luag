package vm

import "github.com/captainlee1024/luag/api"

// loadNil 给指定数量的寄存器置为nil
func loadNil(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()

	a += 1
	// 栈顶推入nil
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		// 置为nil
		vm.Copy(-1, i)
	}
	// 弹出nil
	vm.Pop(1)
}

// 给单个寄存器设置布尔值
// R(A) := (bool)B;
func loadBool(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.PushBoolean(b != 0)
	vm.Replace(a)

	if c != 0 {
		vm.AddPC(1)
	}
}

// 将某个常量加载到寄存器
// R(A) := Kst(Bx)
func loadk(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1

	// 常量加载到栈顶
	// Bx 占18个比特，所以最大只能标识无符号证书262143
	vm.GetConst(bx)
	// 将栈顶的值弹出，替换到寄存器a
	vm.Replace(a)
}

func loadkx(i Instruction, vm api.LuaVM) {
	a, _ := i.ABx()
	a += 1
	// Ax指令中获取ax，来指定常量索引
	// ax 占26个bit，最大可以标识67108864，可以满足大部分情况
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
