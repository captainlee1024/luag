package binchunk

/*
xxd -u -g 1 hw.luac
00000000: 1B 4C 75 61 53 00 19 93 0D 0A 1A 0A 04 08 04 08  .LuaS...........
00000010: 08 78 56 00 00 00 00 00 00 00 00 00 00 00 28 77  .xV...........(w
00000020: 40 01 11 40 68 65 6C 6C 6F 5F 77 6F 72 6C 64 2E  @..@hello_world.
00000030: 6C 75 61 00 00 00 00 00 00 00 00 00 01 02 04 00  lua.............
00000040: 00 00 06 00 40 00 41 40 00 00 24 40 00 01 26 00  ....@.A@..$@..&.
00000050: 80 00 02 00 00 00 04 06 70 72 69 6E 74 04 0E 48  ........print..H
00000060: 65 6C 6C 6F 2C 20 57 6F 72 6C 64 21 01 00 00 00  ello, World!....
00000070: 01 00 00 00 00 00 04 00 00 00 01 00 00 00 01 00  ................
00000080: 00 00 01 00 00 00 01 00 00 00 00 00 00 00 01 00  ................
00000090: 00 00 05 5F 45 4E 56                             ..._ENV
*/

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 祝函数 upvalue 数量
	mainFunc     *Prototype // 主函数原型
}

type header struct {
	// 1B 4C 75 61
	signature [4]byte // 4byte 签名（魔数）十六进制: 0x1B4c7561 go字面量:"\x1bLua"
	// 0x53 (major version * 16 + 3 = 83 -> 16进制 = 0x53)
	version byte // 1byte 版本号 大版本号.小版本号.发布版本号
	// 0x00
	format byte // 1byte 格式号
	// 19 93 0D 0A 1A 0A
	luaData [6]byte // 6byte 0x1993 0x0D 0x0A 0x1A go字面量: "\x19\x93\r\n\x1a\n"
	// TODO:1A 0A?
	// 04
	cintSize byte // 1byte cint数据类型在chunk中占用的字节数
	// 08
	sizetSize byte // 1byte 同上
	// 04
	instrutionSize byte // 1byte 同上
	// 08
	luaIntegerSize byte // 1byte 同上
	// 08
	luaNumberSize byte // 1byte 同上
	// 0x5678
	// 78 56 00 00 00 00 00 00
	luacInt byte // nbyte 接下来的n个字节存放lua整数
	// 370.5
	// 00 00 00 00 00 28 77 40
	luaNum float64 // nbyte 最后n个字节存放lua浮点数
}

type Prototype struct {
	// 11 40 68 65 6C 6C 6F 5F 77 6F 72 6C 64 2E 6C 75 61
	Source string // 源文件名
	// 00 00 00 00
	LineDefined uint32 // 两个cint整型，表示起止行号
	// 00 00 00 00
	LastLineDefined uint32 //
	// 00
	NumParams byte // 1byte 记录函数固定参数个数
	// 01
	IsVararg byte // 1byte 是否是Vararg函数 0:否, 1:是
	// 02
	MaxStackSize byte // 1byte 寄存器数量
	//																					 04 00
	// 00 00 06 00 40 00 41 40 00 00 24 40 00 01 26 00
	// 80 00
	Code []uint32 // 4byte 指令表
	//			 02 00 00 00 04 06 70 72 69 6E 74 04 0E 48
	// 65 6C 6C 6F 2C 20 57 6F 72 6C 64 21
	Constants []interface{} // 常量表
	// 01 00 00 00 01 00
	Upvalues []Upvalue // 2byte Upvalue表
	// 00 00 00 00
	Protos []*Prototype // cint整型表示 子函数原型长度
	// 04 00 00 00
	// 01 00 00 00
	// 01 00 00 00
	// 01 00 00 00
	// 01 00 00 00
	LineInfo []uint32 // 行号表
	// 00 00 00 00
	LocVars []LocVar // 局部变量表
	// 01 00 00 00 05 5F 45 4E 56
	UpvalueNames []string // Upvalue名列表
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string // 变量名
	StartPC uint32 // 起止指令索引
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 校验头部
	reader.readByte()           // 跳过Upvalue数量
	return reader.readProto("") // 读取函数原型
}
