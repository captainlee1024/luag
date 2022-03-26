package state

import "github.com/captainlee1024/luag/api"

type luaStack struct {
	/* virtual stack */
	slots []luaValue
	top   int

	/* call info */
	// pc,protoType 属于函数调用内部状态，所以放在调用帧里合适
	// protoType封装成了closure
	closure *closure
	varargs []luaValue
	pc      int
	openuvs map[int]*upvalue // 记录暂时还处于open状态的upvalue

	/* linked list */
	prev  *luaStack
	state *luaState // 让stack引用state，这样就可以间接访问注册表
}

/*
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}
*/

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

func (stack *luaStack) check(n int) {
	free := len(stack.slots) - stack.top
	for i := free; i < n; i++ {
		stack.slots = append(stack.slots, nil)
	}
}

func (stack *luaStack) push(val luaValue) {
	if stack.top == len(stack.slots) {
		panic("stack overflow!")
	}
	stack.slots[stack.top] = val
	stack.top++
}

func (stack *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			stack.push(vals[i])
		} else {
			stack.push(nil)
		}
	}
}

func (stack *luaStack) pop() luaValue {
	if stack.top < 1 {
		panic("stack underflow!")
	}
	stack.top--
	val := stack.slots[stack.top]
	stack.slots[stack.top] = nil
	return val
}

func (stack *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = stack.pop()
	}
	return vals
}

func (stack *luaStack) absIndex(idx int) int {
	if idx <= api.LUA_REGISTRYINDEX { // 伪索引,用来查找全局变量表的
		return idx
	}

	if idx >= 0 {
		return idx
	}
	return idx + stack.top + 1
}

func (stack *luaStack) isValid(idx int) bool {
	if idx < api.LUA_REGISTRYINDEX { // upbalues
		uvIdx := api.LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		return c != nil && uvIdx < len(c.upvals)
	}

	if idx == api.LUA_REGISTRYINDEX { // 全局注册表
		return true
	}
	absIdx := stack.absIndex(idx)
	return absIdx > 0 && absIdx <= stack.top
}

// get 获取指定寄存器的值
func (stack *luaStack) get(idx int) luaValue {
	if idx < api.LUA_REGISTRYINDEX { // upvalues
		uvIdx := api.LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}

	// 如果是伪索引，返回全局注册表
	if idx == api.LUA_REGISTRYINDEX {
		return stack.state.registry
	}

	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		return stack.slots[absIdx-1]
	}
	return nil
}

func (stack *luaStack) set(idx int, val luaValue) {
	if idx < api.LUA_REGISTRYINDEX { // upvalues
		uvIdx := api.LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}

	if idx == api.LUA_REGISTRYINDEX { // 全局注册表
		stack.state.registry = val.(*luaTable)
		return
	}

	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		stack.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

// 旋转
func (stack *luaStack) reverse(from, to int) {
	slots := stack.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
