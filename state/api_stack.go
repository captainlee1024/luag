package state

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gettop
func (ls *luaState) GetTop() int {
	return ls.stack.top
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_absindex
func (ls *luaState) AbsIndex(idx int) int {
	return ls.stack.absIndex(idx)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_checkstack
func (ls *luaState) CheckStack(n int) bool {
	ls.stack.check(n)
	return true // never fails
}

// [-n, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pop
func (ls *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		ls.stack.pop()
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_copy
func (ls *luaState) Copy(fromIdx, toIdx int) {
	val := ls.stack.get(fromIdx)
	ls.stack.set(toIdx, val)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvalue
func (ls *luaState) PushValue(idx int) {
	val := ls.stack.get(idx)
	ls.stack.push(val)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_replace
func (ls *luaState) Replace(idx int) {
	val := ls.stack.pop()
	ls.stack.set(idx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_insert
func (ls *luaState) Insert(idx int) {
	ls.Rotate(idx, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_remove
func (ls *luaState) Remove(idx int) {
	ls.Rotate(idx, -1)
	ls.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rotate
func (ls *luaState) Rotate(idx, n int) {
	// TODO:
	// 需要旋转的两个部分各自进行reverse翻转
	// 之后整体进行一次reverse翻转
	t := ls.stack.top - 1
	p := ls.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	ls.stack.reverse(p, m)   // 先翻转一次,保证最后整体翻转之后该部分内部保持原来的顺序
	ls.stack.reverse(m+1, t) // 先翻转一次，保证最后整体翻转之后该部分内部保持原来的顺序
	ls.stack.reverse(p, t)   // 两个部分进行整体翻转
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_settop
func (ls *luaState) SetTop(idx int) {
	newTop := ls.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := ls.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			ls.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			ls.stack.push(nil)
		}
	}
}
