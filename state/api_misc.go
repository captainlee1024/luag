package state

func (state *luaState) Len(idx int) {
	val := state.stack.get(idx)
	if s, ok := val.(string); ok {
		state.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", state); ok {
		// 查看值是否有__len元方法，如果有则以值为参数调用元方法，将返回值作为结果
		state.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		state.stack.push(int64(t.len()))
	} else {
		panic("length error!")
	}
}

func (state *luaState) Concat(n int) {
	if n == 0 {
		state.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if state.IsString(-1) && state.IsString(-2) {
				s2 := state.ToString(-1)
				s1 := state.ToString(-2)
				state.stack.pop()
				state.stack.pop()
				state.stack.push(s1 + s2)
				continue
			}

			// 如果不是字符串或者数字，尝试调用__concat元方法，调用规则同二元算术元方法
			b := state.stack.pop()
			a := state.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", state); ok {
				state.stack.push(result)
				continue
			}

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}
