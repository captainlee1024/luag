package api

// 为了不再增加LUA API 方法，我们引入一个新的接口LuaVM
// 用来拓展LuaState接口
type LuaVM interface {
	LuaState
	PC() int          // 返回当前PC 测试时使用
	AddPC(n int)      // 修改PC 用于实现跳转指令
	Fetch() uint32    // 取出当前指令，将PC指向下一条指令
	GetConst(idx int) // 将指定常量推入栈顶
	GetRK(rk int)     // 将指定常量或栈值推入栈顶
}
