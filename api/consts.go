package api

const (
	LUA_MINSTACK = 20
	LUA_MAXSTACK
	LUA_REGISTRYINDEX       = -LUA_MAXSTACK - 100
	LUA_RIDX_GLOBALS  int64 = 2
)

// 与basic types 对应的有typeOf API
/* basic types */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

/* arithmetic functions */
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~

)

/* comparsion functions */
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)
