package nes

import "image/color"

type Instruction struct{
	mnemonic string
	addrMode AddrMode
	cycle    uint64
}

var instructions = map[byte]Instruction{
	0x00: Instruction{"BRK", Implied, 7},
	0x01: Instruction{"ORA", IndirectX, 6},
	0x02: Instruction{"KIL", Accumulator, 2},
	0x04: Instruction{"NOP", Zeropage, 3},
	0x05: Instruction{"ORA", Zeropage, 3},
	0x06: Instruction{"ASL", Zeropage, 5},
	0x08: Instruction{"PHP", Implied, 3},
	0x09: Instruction{"ORA", Immediate, 2},
	0x0A: Instruction{"ASL", Accumulator, 2},
	0x0C: Instruction{"NOP", Absolute, 4},
	0x0D: Instruction{"ORA", Absolute, 4},
	0x0E: Instruction{"ASL", Absolute, 6},
	0x10: Instruction{"BPL", Relative, 2},
	0x11: Instruction{"ORA", IndirectY, 5},
	0x12: Instruction{"KIL", Accumulator, 2},
	0x15: Instruction{"ORA", ZeropageX, 4},
	0x16: Instruction{"ASL", ZeropageX, 6},
	0x18: Instruction{"CLC", Implied, 2},
	0x19: Instruction{"ORA", AbsoluteY, 4},
	0x1D: Instruction{"ORA", AbsoluteX, 4},
	0x1E: Instruction{"ASL", AbsoluteX, 7},
	0x20: Instruction{"JSR", Absolute, 6},
	0x21: Instruction{"AND", IndirectX, 6},
	0x24: Instruction{"BIT", Zeropage, 3},
	0x25: Instruction{"AND", Zeropage, 3},
	0x26: Instruction{"ROL", Zeropage, 5},
	0x28: Instruction{"PLP", Implied, 4},
	0x29: Instruction{"AND", Immediate, 2},
	0x2A: Instruction{"ROL", Accumulator, 2},
	0x2C: Instruction{"BIT", Absolute, 4},
	0x2D: Instruction{"AND", Absolute, 2},
	0x2E: Instruction{"ROL", Absolute, 6},
	0x2F: Instruction{"RLA", Absolute, 6},
	0x30: Instruction{"BMI", Relative, 2},
	0x31: Instruction{"AND", IndirectY, 5},
	0x35: Instruction{"AND", ZeropageX, 4},
	0x36: Instruction{"ROL", ZeropageX, 6},
	0x37: Instruction{"RLA", ZeropageX, 6},
	0x38: Instruction{"SEC", Implied, 2},
	0x39: Instruction{"AND", AbsoluteY, 4},
	0x3D: Instruction{"AND", AbsoluteX, 4},
	0x3E: Instruction{"ROL", AbsoluteX, 7},
	0x40: Instruction{"RTI", Implied, 6},
	0x41: Instruction{"EOR", IndirectX, 6},
	0x43: Instruction{"SRE", IndirectX, 8},
	0x44: Instruction{"NOP", Zeropage, 3},
	0x45: Instruction{"EOR", Zeropage, 3},
	0x46: Instruction{"LSR", Zeropage, 5},
	0x47: Instruction{"SRE", Zeropage, 5},
	0x48: Instruction{"PHA", Implied, 3},
	0x49: Instruction{"EOR", Immediate, 2},
	0x4A: Instruction{"LSR", Accumulator, 2},
	0x4C: Instruction{"JMP", Absolute, 3},
	0x4D: Instruction{"EOR", Absolute, 4},
	0x4E: Instruction{"LSR", Absolute, 6},
	0x50: Instruction{"BVC", Relative, 2},
	0x51: Instruction{"EOR", IndirectY, 5},
	0x55: Instruction{"EOR", ZeropageX, 4},
	0x56: Instruction{"LSR", ZeropageX, 6},
	0x58: Instruction{"CLI", Implied, 2},
	0x59: Instruction{"EOR", AbsoluteY, 4},
	0x5A: Instruction{"NOP", Accumulator, 2},
	0x5D: Instruction{"EOR", AbsoluteX, 4},
	0x5E: Instruction{"LSR", AbsoluteX, 7},
	0x60: Instruction{"RTS", Implied, 6},
	0x61: Instruction{"ADC", IndirectX, 6},
	0x65: Instruction{"ADC", Zeropage, 3},
	0x66: Instruction{"ROR", Zeropage, 5},
	0x68: Instruction{"PLA", Implied, 4},
	0x69: Instruction{"ADC", Immediate, 2},
	0x6A: Instruction{"ROR", Accumulator, 2},
	0x6C: Instruction{"JMP", Indirect, 5},
	0x6D: Instruction{"ADC", Absolute, 4},
	0x6E: Instruction{"ROR", Absolute, 6},
	0x6F: Instruction{"RRA", Absolute, 6},
	0x70: Instruction{"BVS", Relative, 2},
	0x71: Instruction{"ADC", IndirectY, 5},
	0x75: Instruction{"ADC", ZeropageX, 4},
	0x76: Instruction{"ROR", ZeropageX, 6},
	0x78: Instruction{"SEI", Implied, 2},
	0x79: Instruction{"ADC", AbsoluteY, 4},
	0x7D: Instruction{"ADC", AbsoluteX, 4},
	0x7E: Instruction{"ROR", AbsoluteX, 7},
	0x80: Instruction{"NOP", Immediate, 2},
	0x81: Instruction{"STA", IndirectX, 6},
	0x84: Instruction{"STY", Zeropage, 3},
	0x85: Instruction{"STA", Zeropage, 3},
	0x86: Instruction{"STX", Zeropage, 3},
	0x88: Instruction{"DEY", Implied, 2},
	0x8A: Instruction{"TXA", Implied, 2},
	0x8C: Instruction{"STY", Absolute, 4},
	0x8D: Instruction{"STA", Absolute, 4},
	0x8E: Instruction{"STX", Absolute, 4},
	0x90: Instruction{"BCC", Relative, 2},
	0x91: Instruction{"STA", IndirectY, 6},
	0x94: Instruction{"STY", ZeropageX, 4},
	0x95: Instruction{"STA", ZeropageX, 4},
	0x96: Instruction{"STX", ZeropageY, 4},
	0x98: Instruction{"TYA", Implied, 2},
	0x99: Instruction{"STA", AbsoluteY, 5},
	0x9A: Instruction{"TXS", Implied, 2},
	0x9C: Instruction{"SHY", AbsoluteX, 5},
	0x9D: Instruction{"STA", AbsoluteX, 5},
	0xA0: Instruction{"LDY", Immediate, 2},
	0xA1: Instruction{"LDA", IndirectX, 6},
	0xA2: Instruction{"LDX", Immediate, 2},
	0xA3: Instruction{"LAX", IndirectX, 6},
	0xA4: Instruction{"LDY", Zeropage, 3},
	0xA5: Instruction{"LDA", Zeropage, 3},
	0xA6: Instruction{"LDX", Zeropage, 3},
	0xA8: Instruction{"TAY", Implied, 2},
	0xA9: Instruction{"LDA", Immediate, 2},
	0xAA: Instruction{"TAX", Implied, 2},
	0xAC: Instruction{"LDY", Absolute, 4},
	0xAD: Instruction{"LDA", Absolute, 4},
	0xAE: Instruction{"LDX", Absolute, 4},
	0xB0: Instruction{"BCS", Relative, 2},
	0xB1: Instruction{"LDA", IndirectY, 5},
	0xB4: Instruction{"LDY", ZeropageX, 4},
	0xB5: Instruction{"LDA", ZeropageX, 4},
	0xB6: Instruction{"LDX", ZeropageY, 4},
	0xB8: Instruction{"CLV", Implied, 2},
	0xB9: Instruction{"LDA", AbsoluteY, 4},
	0xBA: Instruction{"TSX", Implied, 2},
	0xBC: Instruction{"LDY", AbsoluteX, 4},
	0xBD: Instruction{"LDA", AbsoluteX, 4},
	0xBE: Instruction{"LDX", AbsoluteY, 4},
	0xBF: Instruction{"LAX", AbsoluteY, 4},
	0xC0: Instruction{"CPY", Immediate, 2},
	0xC1: Instruction{"CMP", IndirectX, 6},
	0xC4: Instruction{"CPY", Zeropage, 3},
	0xC5: Instruction{"CMP", Zeropage, 3},
	0xC6: Instruction{"DEC", Zeropage, 5},
	0xC8: Instruction{"INY", Implied, 2},
	0xC9: Instruction{"CMP", Immediate, 2},
	0xCA: Instruction{"DEX", Implied, 2},
	0xCC: Instruction{"CPY", Absolute, 4},
	0xCD: Instruction{"CMP", Absolute, 4},
	0xCE: Instruction{"DEC", Absolute, 6},
	0xD0: Instruction{"BNE", Relative, 2},
	0xD1: Instruction{"CMP", IndirectY, 5},
	0xD3: Instruction{"DCP", IndirectY, 8},
	0xD5: Instruction{"CMP", ZeropageX, 4},
	0xD6: Instruction{"DEC", ZeropageX, 6},
	0xD8: Instruction{"CLD", Implied, 2},
	0xD9: Instruction{"CMP", AbsoluteY, 4},
	0xDD: Instruction{"CMP", AbsoluteX, 4},
	0xDE: Instruction{"DEC", AbsoluteX, 7},
	0xE0: Instruction{"CPX", Immediate, 2},
	0xE1: Instruction{"SBC", IndirectX, 6},
	0xE4: Instruction{"CPX", Zeropage, 3},
	0xE5: Instruction{"SBC", Zeropage, 3},
	0xE6: Instruction{"INC", Zeropage, 5},
	0xE8: Instruction{"INX", Implied, 2},
	0xE9: Instruction{"SBC", Immediate, 2},
	0xEA: Instruction{"NOP", Implied, 2},
	0xEB: Instruction{"SBC", Immediate, 2},
	0xEC: Instruction{"CPX", Absolute, 4},
	0xED: Instruction{"SBC", Absolute, 4},
	0xEE: Instruction{"INC", Absolute, 6},
	0xF0: Instruction{"BEQ", Relative, 2},
	0xF1: Instruction{"SBC", IndirectY, 5},
	0xF5: Instruction{"SBC", ZeropageX, 4},
	0xF6: Instruction{"INC", ZeropageX, 6},
	0xF8: Instruction{"SED", Implied, 2},
	0xF9: Instruction{"SBC", AbsoluteY, 4},
	0xFD: Instruction{"SBC", AbsoluteX, 4},
	0xFE: Instruction{"INC", AbsoluteX, 7},
	0xFF: Instruction{"ISC", AbsoluteX, 7},
}

var SystemPalette [64]color.RGBA

func initSystemPallete() {
	colors := []uint32{
		0x666666, 0x002A88, 0x1412A7, 0x3B00A4, 0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
		0x333500, 0x0B4800, 0x005200, 0x004F08, 0x00404D, 0x000000, 0x000000, 0x000000,
		0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE, 0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
		0x6B6D00, 0x388700, 0x0C9300, 0x008F32, 0x007C8D, 0x000000, 0x000000, 0x000000,
		0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF, 0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
		0xBCBE00, 0x88D800, 0x5CE430, 0x45E082, 0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
		0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF, 0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
		0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC, 0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
	}

	for i, c := range colors {
		r := byte(c >> 16)
		g := byte(c >> 8)
		b := byte(c)
		SystemPalette[i] = color.RGBA{R: r, G: g, B: b, A: 0xFF}
	}
}