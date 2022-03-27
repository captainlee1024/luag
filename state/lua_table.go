package state

import (
	"math"

	"github.com/captainlee1024/luag/number"
)

type luaTable struct {
	metatable *luaTable // 元表
	arr       []luaValue
	_map      map[luaValue]luaValue

	/* used by next() */
	keys    map[luaValue]luaValue
	lastKey luaValue
	changed bool
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, nArr)
	}

	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}

	return t
}

func (lt *luaTable) hasMetafield(fieldName string) bool {
	return lt.metatable != nil &&
		lt.metatable.get(fieldName) != nil
}

func (lt *luaTable) len() int {
	return len(lt.arr)
}

func (lt *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(lt.arr)) {
			return lt.arr[idx-1]
		}
	}

	return lt._map[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}

	return key
}

func (lt *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("lt index is nil!")
	}

	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("lt index is NaN!")
	}

	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(lt.arr))
		if idx <= arrLen {
			lt.arr[idx-1] = val
			if idx == arrLen && val == nil {
				lt._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(lt._map, key)
			if val != nil {
				lt.arr = append(lt.arr, val)
				lt._expandArray()
			}
			return
		}
	}

	if val != nil {
		if lt._map == nil {
			lt._map = make(map[luaValue]luaValue, 8)
		}
		lt._map[key] = val
	} else {
		delete(lt._map, key)
	}
}

func (lt *luaTable) _shrinkArray() {
	for i := len(lt.arr) - 1; i >= 0; i-- {
		if lt.arr[i] == nil {
			lt.arr[i] = lt.arr[0:i]
		} else {
			break
		}
	}
}

func (lt *luaTable) _expandArray() {
	for idx := int64(len(lt.arr)) + 1; true; idx++ {
		if val, found := lt._map[idx]; found {
			delete(lt._map, idx)
			lt.arr = append(lt.arr, val)
		} else {
			break
		}
	}
}

func (lt *luaTable) nextKey(key luaValue) luaValue {
	if lt.keys == nil || (key == nil && lt.changed) {
		lt.initKeys()
		lt.changed = false
	}

	nextKey := lt.keys[key]
	if nextKey == nil && key != nil && key != lt.lastKey {
		panic("invalid key to 'next'")
	}

	return nextKey
}

func (lt *luaTable) initKeys() {
	lt.keys = make(map[luaValue]luaValue)
	var key luaValue = nil
	for i, v := range lt.arr {
		if v != nil {
			lt.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}
	for k, v := range lt._map {
		if v != nil {
			lt.keys[key] = k
			key = k
		}
	}
	lt.lastKey = key
}
