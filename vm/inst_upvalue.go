package vm

import "github.com/captainlee1024/luag/api"

/*
func getTabUp(i Instruction, vm api.LuaVM) {
	a, _, c := i.ABC()
	a += 1

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
}
*/

/*
func getUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(api.LuaUpvalueIndex(b), a)
}

func setUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(a, api.LuaUpvalueIndex(b))
}

func getTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(api.LuaUpvalueIndex(b))
	vm.Replace(a)
}

func setTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(api.LuaUpvalueIndex(a))
}
*/

// R(A) := UpValue[B]
func getUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(api.LuaUpvalueIndex(b), a)
}

// UpValue[B] := R(A)
func setUpval(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, api.LuaUpvalueIndex(b))
}

// R(A) := UpValue[B][RK(C)]
func getTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(api.LuaUpvalueIndex(b))
	vm.Replace(a)
}

// UpValue[A][RK(B)] := RK(C)
func setTabUp(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(api.LuaUpvalueIndex(a))
}
