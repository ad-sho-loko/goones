package nes

import (
	"fmt"
)

type Cpu struct{
	A     byte
	X     byte
	Y     byte
	S     byte
	P     byte
	PC    word
	cycle uint64
	bus   *Bus
	intrrupt func()
}

const(
	Carry    = 0x01
	Zero     = 0x02
	Irq      = 0x04
	Decimal  = 0x08
	Break    = 0x10
	Reserved = 0x20
	Overflow = 0x40
	Negative = 0x80
)

type AddrMode uint

const(
	Accumulator AddrMode = iota
	Immediate
	Zeropage
	ZeropageX
	ZeropageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	Implied
	Relative
)

func (a AddrMode) String() string{
	switch a {
	case Accumulator: return "Accumulator"
	case Immediate: return "Immediate"
	case Zeropage: return "Zeropage"
	case ZeropageX : return "ZeropageX"
	case ZeropageY : return "ZeropageY"
	case Absolute : return "Absolute"
	case AbsoluteX : return "AbsoluteX"
	case AbsoluteY : return "AbsoluteY"
	case Indirect : return "Indirect"
	case IndirectX : return "IndirectX"
	case IndirectY : return "IndirectY"
	case Implied : return "Implied"
	case Relative : return "Relative"
	}
	panic("Unable to reach here")
}

func NewCpu(bus *Bus) *Cpu{
	return &Cpu{
		bus:bus,
	}
}

func (c *Cpu) status(kind byte) byte{
	if c.P & kind != 0{
		return 0x01
	}else{
		return 0x00
	}
}

func (c *Cpu) setBit(kind byte){
	c.P |= kind
}

func (c *Cpu) unsetBit(kind byte){
	c.P &= 0xFF - kind
}

func (c *Cpu) isCarry() bool{
	return c.P & Carry != 0
}

func (c *Cpu) isZero() bool{
	return c.P & Zero != 0
}

func (c *Cpu) isNegative() bool{
	return c.P & Negative != 0
}

func (c *Cpu) isOverflow() bool{
	return c.P & Overflow != 0
}

func (c *Cpu) isIrqForbitten() bool{
	return c.P & Irq != 0
}

func (c *Cpu) updateNZ(b byte){
	if b == 0x00 {
		c.setBit(Zero)
	}else{
		c.unsetBit(Zero)
	}

	if b & 0x80 != 0 {
		c.setBit(Negative)
	}else{
		c.unsetBit(Negative)
	}
}

func (c *Cpu) lda(w word){
	c.A = c.bus.Load(w)
	c.updateNZ(c.A)
}

func (c *Cpu) ldx(w word){
	c.X = c.bus.Load(w)
	c.updateNZ(c.X)
}

func (c *Cpu) ldy(w word){
	c.Y = c.bus.Load(w)
	c.updateNZ(c.Y)
}

func (c *Cpu) sta(addr word){
	c.bus.Store(addr, c.A)
}

func (c *Cpu) stx(addr word){
	c.bus.Store(addr, c.X)
}

func (c *Cpu) sty(addr word){
	c.bus.Store(addr, c.Y)
}

func (c *Cpu) tax(){
	c.X = c.A
	c.updateNZ(c.X)
}

func (c *Cpu) tay(){
	c.Y = c.A
	c.updateNZ(c.Y)
}

func (c *Cpu) tsx(){
	c.X = c.S
	c.updateNZ(c.X)
}

func (c *Cpu) txa(){
	c.A = c.X
	c.updateNZ(c.A)
}

func (c *Cpu) txs(){
	c.S = c.X
	// c.updateNZ(c.S)
	// nothing updates
}

func (c *Cpu) tya(){
	c.A = c.Y
	c.updateNZ(c.A)
}

func (c *Cpu) adc(w word){
	a := c.A
	b := c.bus.Load(w)
	cy := c.status(Carry)
	c.A = a + b + cy
	c.updateNZ(c.A)

	if int(a)+int(b)+int(cy) > 0xFF{
		c.setBit(Carry)
	}else{
		c.unsetBit(Carry)
	}

	if (a^b) & 0x80 == 0 && (a^c.A)&0x80 != 0{
		c.setBit(Overflow)
	}else{
		c.unsetBit(Overflow)
	}

}

func (c *Cpu) and(w word){
	c.A = c.A & c.bus.Load(w)
	c.updateNZ(c.A)
}

func (c *Cpu) asl(isAccumulator bool, addr word){
	if isAccumulator{
		if c.A >> 7 & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		c.A <<= 1
		c.updateNZ(c.A)

	}else{
		v := c.bus.Load(addr)

		if v >> 7 & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		v <<= 1
		c.bus.Store(addr, v)
		c.updateNZ(v)
	}
}

func (c *Cpu) bit(addr word){
	// 特殊なレジスタ操作が必要なのでロジックを個別化する
	v := c.bus.Load(addr)

	if (v >> 6) & 0x01 != 0{
		c.setBit(Overflow)
	}else{
		c.unsetBit(Overflow)
	}

	if (v & c.A) == 0x00{
		c.setBit(Zero)
	}else{
		c.unsetBit(Zero)
	}

	if (v & 0x80) != 0{
		c.setBit(Negative)
	}else{
		c.unsetBit(Negative)
	}

}

func (c *Cpu) compare(left byte, right byte){
	c.updateNZ(left - right)
	if left >= right{
		c.setBit(Carry)
	} else{
		c.unsetBit(Carry)
	}
}

func (c *Cpu) cmp(addr word){
	v := c.bus.Load(addr)
	c.compare(c.A, v)
}

func (c *Cpu) cpx(addr word){
	v := c.bus.Load(addr)
	c.compare(c.X, v)
}

func (c *Cpu) cpy(addr word){
	v := c.bus.Load(addr)
	c.compare(c.Y, v)
}

func (c *Cpu) dec(addr word){
	v := c.bus.Load(addr)
	c.bus.Store(addr, v-1)
	c.updateNZ(v-1)
}

func (c *Cpu) dex(){
	c.X--
	c.updateNZ(c.X)
}

func (c *Cpu) dey(){
	c.Y--
	c.updateNZ(c.Y)
}

func (c *Cpu) eor(w word){
	c.A ^= c.bus.Load(w)
	c.updateNZ(c.A)
}

func (c *Cpu) inc(addr word){
	v := c.bus.Load(addr)
	c.bus.Store(addr, v+1)
	c.updateNZ(v+1)
}

func (c *Cpu) inx(){
	c.X++
	c.updateNZ(c.X)
}

func (c *Cpu) iny(){
	c.Y++
	c.updateNZ(c.Y)
}

func (c *Cpu) lsr(isAccumulator bool, addr word){
	if isAccumulator{
		if c.A & 1 == 1{
			c.setBit(Carry)
		} else {
			c.unsetBit(Carry)
		}

		c.A >>= 1
		c.updateNZ(c.A)

	}else{
		v := c.bus.Load(addr)
		if v & 1 == 1{
			c.setBit(Carry)
		} else {
			c.unsetBit(Carry)
		}
		v>>=1
		c.bus.Store(addr, v)
		c.updateNZ(v)
	}
}

func (c *Cpu) ora(w word){
	c.A = c.A | c.bus.Load(w)
	c.updateNZ(c.A)
}

func (c *Cpu) rol(isAccumulator bool, addr word){
	if isAccumulator {
		cv := c.status(Carry)

		if (c.A >> 7) & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		c.A = (c.A << 1) | cv
		c.updateNZ(c.A)

	} else {
		cv := c.status(Carry)
		value := c.bus.Load(addr)

		if (value >> 7) & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		value = (value << 1) | cv
		c.bus.Store(addr, value)
		c.updateNZ(value)
	}
}

func (c *Cpu) ror(isAccumulator bool, addr word) {
	if isAccumulator {
		cv := c.status(Carry)

		if c.A & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		c.A = (c.A >> 1) | (cv << 7)
		c.updateNZ(c.A)
	} else {
		cv := c.status(Carry)
		value := c.bus.Load(addr)

		if value & 1 == 1{
			c.setBit(Carry)
		}else{
			c.unsetBit(Carry)
		}

		value = (value >> 1) | (cv << 7)
		c.bus.Store(addr, value)
		c.updateNZ(value)
	}
}

func (c *Cpu) sbc(addr word){
	a := c.A
	b := c.bus.Load(addr)
	cy := c.status(Carry)
	c.A = a - b - (1 - cy)

	if int(a)-int(b)-int(1-cy) >= 0{
		c.setBit(Carry)
	} else {
		c.unsetBit(Carry)
	}

	if (a^b)&0x80 != 0 && (a^c.A)&0x80 != 0{
		c.setBit(Overflow)
	}else{
		c.unsetBit(Overflow)
	}

	c.updateNZ(c.A)
}

func (c *Cpu) push(b byte){
	c.bus.Store(0x100 | word(c.S), b)
	c.S--
}

func (c *Cpu) pushWord(w word){
	h := byte(w >> 8)
	l := byte(w & 0xFF)
	c.push(h)
	c.push(l)
}

func (c *Cpu) pop() byte{
	c.S++
	return c.bus.Load(0x100 | word(c.S))
}

func (c *Cpu) popWord() word {
	l := word(c.pop())
	h := word(c.pop())
	return h << 8 | l
}

func (c *Cpu) pha(){
	c.push(c.A)
}

func (c *Cpu) php(){
	// sets 0x10 always
	c.push(c.P | 0x10)
}

func (c *Cpu) pla(){
	c.A = c.pop()
	c.updateNZ(c.A)
}

func (c *Cpu) plp(){
	c.P = c.pop() & 0xEF | 0x20
}

func (c *Cpu) jmp(addr word){
	c.PC = addr
}

func (c *Cpu) jsr(addr word){
	c.pushWord(c.PC - 1)
	c.jmp(addr)
}

func (c *Cpu) rts(){
	c.PC = c.popWord() + 1
}

func (c *Cpu) rti() {
	c.P = c.pop() & 0xEF | 0x20
	c.PC = c.popWord()
}

func (c *Cpu) branch(w word){
	c.PC = w
}

func (c *Cpu) bcs(w word){
	if c.isCarry(){
		c.branch(w)
	}
}

func (c *Cpu) bcc(w word){
	if !c.isCarry(){
		c.branch(w)
	}
}

func (c *Cpu) beq(w word){
	if c.isZero(){
		c.branch(w)
	}
}

func (c *Cpu) bne(w word){
	if !c.isZero(){
		c.branch(w)
	}
}

func (c *Cpu) bmi(w word){
	if c.isNegative(){
		c.branch(w)
	}
}

func (c *Cpu) bpl(w word){
	if !c.isNegative(){
		c.branch(w)
	}
}

func (c *Cpu) bvc(w word){
	if !c.isOverflow(){
		c.branch(w)
	}
}

func (c *Cpu) bvs(w word){
	if c.isOverflow(){
		c.branch(w)
	}
}

func (c *Cpu) clc(){
	c.unsetBit(Carry)
}

func (c *Cpu) cld(){
	c.unsetBit(Decimal)
}

func (c *Cpu) cli(){
	c.unsetBit(Irq)
}

func (c *Cpu) clv(){
	c.unsetBit(Overflow)
}

func (c *Cpu) sec(){
	c.setBit(Carry)
}

func (c *Cpu) sed(){
	c.setBit(Decimal)
}

func (c *Cpu) sei(){
	c.setBit(Irq)
}

func (c *Cpu) rra(){
	// nop
}

func (c *Cpu) sre(){
	// nop
}

func (c *Cpu) dcp(){
	// nop
}

func (c *Cpu) rla(){
	// nop
}

func (c *Cpu) shy(){
	// nop
}

func (c *Cpu) lax(){
	// nop
}

func (c *Cpu) kil(){
	// nop
}

func (c *Cpu) isc(){
	// nop
}

func (c *Cpu) sax(){
	// nop
}

func (c *Cpu) las(){
	// nop
}

func (c *Cpu) slo(){
	// nop
}

func (c *Cpu) anc(){
	// nop
}

func (c *Cpu) brk(){
	c.pushWord(c.PC)
	c.php()
	c.sei()
	addr := c.bus.Loadw(0xFFFE)
	c.jmp(addr)
	// c.setBit(Irq)
	// c.setBit(Break)
}

func (c *Cpu) nop(){
	// nop
}

func (c *Cpu) nmi(){
	c.pushWord(c.PC)
	c.php()
	addr := c.bus.Loadw(0xFFFA)
	c.jmp(addr)
	c.setBit(Irq)
	// c.unsetBit(Break)
	c.cycle+=7
}

func (c *Cpu) reset(){
	c.pushWord(c.PC)
	c.php()
	addr := c.bus.Loadw(0xFFFC)
	c.jmp(addr)
	c.setBit(Irq)
	// c.unsetBit(Break)
	c.cycle += 7
}

func (c *Cpu) irq(){
	c.pushWord(c.PC)
	c.php()
	addr := c.bus.Loadw(0xFFFE)
	c.jmp(addr)
	c.setBit(Irq)
	c.cycle += 7
}

func (c *Cpu) InterruptNmi(){
	c.intrrupt = c.nmi
}

func (c *Cpu) interruptReset(){
	if !c.isIrqForbitten(){
		c.intrrupt = c.reset
	}
}

func (c *Cpu) interruptIrq(){
	if !c.isIrqForbitten(){
		c.intrrupt = c.irq
	}
}

func (c *Cpu) decode(b byte) Instruction{
	i := instructions[b]
	if i.mnemonic == ""{
		abort("panic: unknown `%x` was decoded.", b)
	}
	return i
}

func (c *Cpu) advance(mode AddrMode){
	switch mode {
	case Accumulator, Implied:
		c.PC += 1
	case Immediate, Zeropage, ZeropageX, ZeropageY, Relative, IndirectX, IndirectY:
		c.PC += 2
	case Absolute, AbsoluteX, AbsoluteY, Indirect:
		c.PC += 3
	default:
		abort("panic: unknown addrMode `%s` was called when advance", mode)
	}
}

func (c *Cpu) solveAddrMode(mode AddrMode) word {
	switch mode {
	case Accumulator:
		return 0x00
	case Implied:
		return 0x00
	case Immediate:
		return c.PC + 1
	case Relative:
		offset := word(c.bus.Load(c.PC + 1))
		if offset < 0x80 {
			return c.PC + 2 + offset
		} else {
			return c.PC + 2 + offset - 0x100
		}
	case Zeropage:
		return word(c.bus.Load(c.PC + 1))
	case ZeropageX:
		return word(c.bus.Load(c.PC + 1) + c.X) & 0xFF
	case ZeropageY:
		return word(c.bus.Load(c.PC + 1) + c.Y) & 0xFF
	case Absolute:
		return c.bus.Loadw(c.PC + 1)
	case AbsoluteX:
		return c.bus.Loadw(c.PC + 1) + word(c.X)
	case AbsoluteY:
		return c.bus.Loadw(c.PC + 1) + word(c.Y)
	case Indirect:
		return c.bus.BugLoadw(c.bus.Loadw(c.PC + 1))
	case IndirectX:
		return c.bus.BugLoadw(word(c.bus.Load(c.PC + 1) + c.X))
	case IndirectY:
		return c.bus.BugLoadw(word(c.bus.Load(c.PC + 1))) + word(c.Y)
	default:
		abort("panic: unknown addrMode `%s` was called when solving", mode)
	}
	panic("Unable to reach here")
}

func (c *Cpu) execute(inst Instruction, w word){
	switch inst.mnemonic {
	case "LDA":
		c.lda(w)
	case "LDX":
		c.ldx(w)
	case "LDY":
		c.ldy(w)
	case "STA":
		c.sta(w)
	case "STX":
		c.stx(w)
	case "STY":
		c.sty(w)
	case "CPX":
		c.cpx(w)
	case "TXS":
		c.txs()
	case "TSX":
		c.tsx()
	case "TYA":
		c.tya()
	case "TAX":
		c.tax()
	case "TXA":
		c.txa()
	case "TAY":
		c.tay()
	case "BIT":
		c.bit(w)
	case "ADC":
		c.adc(w)
	case "SBC":
		c.sbc(w)
	case "AND":
		c.and(w)
	case "ORA":
		c.ora(w)
	case "EOR":
		c.eor(w)
	case "INC":
		c.inc(w)
	case "INX":
		c.inx()
	case "INY":
		c.iny()
	case "DEC":
		c.dec(w)
	case "DEX":
		c.dex()
	case "DEY":
		c.dey()
	case "CMP":
		c.cmp(w)
	case "ASL":
		c.asl(inst.addrMode == Accumulator, w)
	case "ROR":
		c.ror(inst.addrMode == Accumulator, w)
	case "ROL":
		c.rol(inst.addrMode == Accumulator, w)
	case "LSR":
		c.lsr(inst.addrMode == Accumulator, w)
	case "CPY":
		c.cpy(w)
	case "CLC":
		c.clc()
	case "CLD":
		c.cld()
	case "CLI":
		c.cli()
	case "CLV":
		c.clv()
	case "SEC":
		c.sec()
	case "SED":
		c.sed()
	case "SEI":
		c.sei()
	// ジャンプ命令
	case "JMP":
		c.jmp(w)
	case "JSR":
		c.jsr(w)
	case "RTS":
		c.rts()
	case "RTI":
		c.rti()
	case "PLA":
		c.pla()
	case "PHA":
		c.pha()
	case "PHP":
		c.php()
	case "PLP":
		c.plp()
	case "BCC":
		c.bcc(w)
	case "BCS":
		c.bcs(w)
	case "BEQ":
		c.beq(w)
	case "BMI":
		c.bmi(w)
	case "BNE":
		c.bne(w)
	case "BPL":
		c.bpl(w)
	case "BVC":
		c.bvc(w)
	case "BVS":
		c.bvs(w)
	case "BRK":
		c.brk()
	case "NOP":
		c.nop()
	case "RRA":
		c.rra() // do nothing
	case "SRE":
		c.sre() // do nothing
	case "DCP":
		c.dcp() // do nothing
	case "RLA":
		c.rla() // do nothing
	case "SHY":
		c.shy() // do nothing
	case "LAX":
		c.lax() // do nothing
	case "KIL":
		c.lax() // do nothing
	case "ISC":
		c.isc() // do nothing
	case "SAX":
		c.sax() // do nothing
	case "SLO":
		c.sax() // do nothing
	case "LAS":
		c.las() // do nothing
	case "ANC":
		c.anc() // do nothing
	default:
		abort("panic: unknown mnemonic `%s` was invoked.", inst.mnemonic)
	}
}

func (c *Cpu) dump(b byte, arg word, mne string, mode AddrMode){
	fmt.Printf("[PC:0x%4x, A:0x%2x, X:0x%2x, Y:0x%2x SP:0x%x P:0x%x CYC:%d] " +
		"%2x, %4x ## %s (%s) \n",
		c.PC, c.A, c.X, c.Y, c.S, c.P, c.cycle, b, arg, mne, mode)
}
