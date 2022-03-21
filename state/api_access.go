package state

import (
	"fmt"

	"github.com/captainlee1024/luag/api"
)

/*与push类方法相反，access方法需要使用栈里面的数据*/

// TODO

func (state *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (state *luaState) Type(idx int) api.LuaType {
	// 判断索引有效性
	if state.stack.isValid(idx) {
		val := state.stack.get(idx)
		return typeOf(val)
	}

	return api.LUA_TNONE
}

func (state *luaState) IsNone(idx int) bool {
	// 根据索引获取栈内数据，然后再判断类型
	return state.Type(idx) == api.LUA_TNONE
}

func (state *luaState) IsNil(idx int) bool {
	return state.Type(idx) == api.LUA_TNIL
}

func (state *luaState) IsNoneOrNil(idx int) bool {
	return state.Type(idx) <= api.LUA_TNIL
}

func (state *luaState) IsBoolean(idx int) bool {
	return state.Type(idx) == api.LUA_TBOOLEAN
}

func (state *luaState) IsInteger(idx int) bool {
	val := state.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (state *luaState) IsNumber(idx int) bool {
	return state.Type(idx) == api.LUA_TNUMBER
}

func (state *luaState) IsString(idx int) bool {
	t := state.Type(idx)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

func (state *luaState) IsTable(idx int) bool {
	return state.Type(idx) == api.LUA_TTABLE
}

func (state *luaState) IsThread(idx int) bool {
	return state.Type(idx) == api.LUA_TTHREAD
}

func (state *luaState) IsFunction(idx int) bool {
	return state.Type(idx) == api.LUA_TFUNCTION
}

func (state *luaState) ToBoolean(idx int) bool {
	val := state.stack.get(idx)
	return converToBoolean(val)
}

func (state *luaState) ToInteger(idx int) int64 {
	i, _ := state.ToIntegerX(idx)
	return i
}

func (state *luaState) ToIntegerX(idx int) (int64, bool) {
	val := state.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

func (state *luaState) ToNumber(idx int) float64 {
	n, _ := state.ToNumberX(idx)
	return n
}

func (state *luaState) ToNumberX(idx int) (float64, bool) {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (state *luaState) ToString(idx int) string {
	s, _ := state.ToStringX(idx)
	return s
}

func (state *luaState) ToStringX(idx int) (string, bool) {
	val := state.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		state.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
