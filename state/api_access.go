package state

import (
	"fmt"

	"github.com/captainlee1024/luag/api"
)

/*与push类方法相反，access方法需要使用栈里面的数据*/

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawlen
func (ls *luaState) RawLen(idx int) uint {
	val := ls.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

// TODO

func (ls *luaState) TypeName(tp api.LuaType) string {
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

func (ls *luaState) Type(idx int) api.LuaType {
	// 判断索引有效性
	if ls.stack.isValid(idx) {
		val := ls.stack.get(idx)
		return typeOf(val)
	}

	return api.LUA_TNONE
}

func (ls *luaState) IsNone(idx int) bool {
	// 根据索引获取栈内数据，然后再判断类型
	return ls.Type(idx) == api.LUA_TNONE
}

func (ls *luaState) IsNil(idx int) bool {
	return ls.Type(idx) == api.LUA_TNIL
}

func (ls *luaState) IsNoneOrNil(idx int) bool {
	return ls.Type(idx) <= api.LUA_TNIL
}

func (ls *luaState) IsBoolean(idx int) bool {
	return ls.Type(idx) == api.LUA_TBOOLEAN
}

func (ls *luaState) IsInteger(idx int) bool {
	val := ls.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (ls *luaState) IsNumber(idx int) bool {
	return ls.Type(idx) == api.LUA_TNUMBER
}

func (ls *luaState) IsString(idx int) bool {
	t := ls.Type(idx)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

func (ls *luaState) IsTable(idx int) bool {
	return ls.Type(idx) == api.LUA_TTABLE
}

func (ls *luaState) IsThread(idx int) bool {
	return ls.Type(idx) == api.LUA_TTHREAD
}

func (ls *luaState) IsGoFunction(idx int) bool {
	val := ls.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

func (ls *luaState) IsFunction(idx int) bool {
	return ls.Type(idx) == api.LUA_TFUNCTION
}

func (ls *luaState) ToBoolean(idx int) bool {
	val := ls.stack.get(idx)
	return convertToBoolean(val)
}

func (ls *luaState) ToInteger(idx int) int64 {
	i, _ := ls.ToIntegerX(idx)
	return i
}

func (ls *luaState) ToIntegerX(idx int) (int64, bool) {
	val := ls.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

func (ls *luaState) ToNumber(idx int) float64 {
	n, _ := ls.ToNumberX(idx)
	return n
}

func (ls *luaState) ToNumberX(idx int) (float64, bool) {
	val := ls.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (ls *luaState) ToString(idx int) string {
	s, _ := ls.ToStringX(idx)
	return s
}

func (ls *luaState) ToStringX(idx int) (string, bool) {
	val := ls.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		ls.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (ls *luaState) ToGoFunction(idx int) api.GoFunction {
	val := ls.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}
