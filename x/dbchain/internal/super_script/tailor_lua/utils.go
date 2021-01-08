package tailor_lua

const opSizeA = 8
const opSizeB = 9
const opSizeC = 9
const opSizeBx = 18

const opMaxArgsA = (1 << opSizeA) - 1
const opMaxArgsC = (1 << opSizeC) - 1
const opMaxArgBx = (1 << opSizeBx) - 1
const opMaxArgSbx = opMaxArgBx >> 1
const opBitRk = 1 << (opSizeB - 1)
const opMaxIndexRk = opBitRk - 1

func opGetArgC(inst uint32) int {
	return int(inst>>9) & 0x1ff
}

func opGetArgB(inst uint32) int {
	return int(inst & 0x1ff)
}

func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func opCreateABC(op int, a int, b int, c int) uint32 {
	var inst uint32 = 0
	opSetOpCode(&inst, op)
	opSetArgA(&inst, a)
	opSetArgB(&inst, b)
	opSetArgC(&inst, c)
	return inst
}

func opCreateABx(op int, a int, bx int) uint32 {
	var inst uint32 = 0
	opSetOpCode(&inst, op)
	opSetArgA(&inst, a)
	opSetArgBx(&inst, bx)
	return inst
}

func opCreateASbx(op int, a int, sbx int) uint32 {
	var inst uint32 = 0
	opSetOpCode(&inst, op)
	opSetArgA(&inst, a)
	opSetArgSbx(&inst, sbx)
	return inst
}

func opSetOpCode(inst *uint32, opcode int) {
	*inst = (*inst & 0x3ffffff) | uint32(opcode<<26)
}

func opSetArgA(inst *uint32, arg int) {
	*inst = (*inst & 0xfc03ffff) | uint32((arg&0xff)<<18)
}

func opSetArgBx(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc0000) | uint32(arg&0x3ffff)
}

func opSetArgSbx(inst *uint32, arg int) {
	opSetArgBx(inst, arg+opMaxArgSbx)
}

func opSetArgB(inst *uint32, arg int) {
	*inst = (*inst & 0xfffffe00) | uint32(arg&0x1ff)
}

func opSetArgC(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc01ff) | uint32((arg&0x1ff)<<9)
}

func opIsK(value int) bool {
	return bool((value & opBitRk) != 0)
}

func opRkAsk(value int) int {
	return value | opBitRk
}

func opGetOpCode(inst uint32) int {
	return int(inst >> 26)
}

func opGetArgSbx(inst uint32) int {
	return opGetArgBx(inst) - opMaxArgSbx
}

func opGetArgBx(inst uint32) int {
	return int(inst & 0x3ffff)
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func opGetArgA(inst uint32) int {
	return int(inst>>18) & 0xff
}
