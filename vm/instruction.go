package vm

import "github.com/captainlee1024/luag/api"

type Instruction uint32

// Opcode 从指令中提取操作码
func (instr Instruction) Opcode() int {
	return int(instr & 0x3F)
}

const MAXARG_Bx = 1<<18 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/

// ABC 从iABC指令提取参数
func (instr Instruction) ABC() (a, b, c int) {
	a = int(instr >> 6 & 0xFF)
	c = int(instr >> 14 & 0x1FF)
	b = int(instr >> 23 & 0x1FF)
	return
}

// ABx 从iABx指令提取参数
func (instr Instruction) ABx() (a, bx int) {
	a = int(instr >> 6 & 0xFF)
	bx = int(instr >> 14)
	return
}

// AsBx 从iAsBx指令提取参数
func (instr Instruction) AsBx() (a, sbx int) {
	a, bx := instr.ABx()
	return a, bx - MAXARG_sBx
}

// Ax 从iAx指令提取参数
func (instr Instruction) Ax() int {
	return int(instr >> 6)
}

func (instr Instruction) OpName() string {
	return opcodes[instr.Opcode()].name
}

func (instr Instruction) OpMode() byte {
	return opcodes[instr.Opcode()].opMode
}

func (instr Instruction) BMode() byte {
	return opcodes[instr.Opcode()].argBMode
}

func (instr Instruction) CMode() byte {
	return opcodes[instr.Opcode()].argCMode
}

func (instr Instruction) Execute(vm api.LuaVM) {
	action := opcodes[instr.Opcode()].action
	if action != nil {
		action(instr, vm)
	} else {
		panic(instr.OpName())
	}
}
