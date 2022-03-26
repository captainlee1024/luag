package state

import (
	"github.com/captainlee1024/luag/api"
	"github.com/captainlee1024/luag/binchunk"
)

type closure struct {
	proto  *binchunk.Prototype
	goFunc api.GoFunction
	upvals []*upvalue // upvalue
}

type upvalue struct {
	val *luaValue
}

/*
func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f api.GoFunction) *closure {
	return &closure{goFunc: f}
}
*/

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f api.GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
