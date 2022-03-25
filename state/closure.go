package state

import (
	"github.com/captainlee1024/luag/api"
	"github.com/captainlee1024/luag/binchunk"
)

type closure struct {
	proto  *binchunk.Prototype
	goFunc api.GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f api.GoFunction) *closure {
	return &closure{goFunc: f}
}
