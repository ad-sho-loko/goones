package nes

import "fmt"

type Cpu struct{
	A   byte
	X   byte
	Y   byte
	S   byte
	P   byte
	PC  Word
	cycle uint64
	bus *Bus
}

const(
	Carry    = 0x01
	Zero     = 0x02
	Irq      = 0x04
	Decimal  = 0x08
	Braak    = 0x10
	Reserved = 0x20
	Overflow = 0x40
	Negative = 0x80
)

type AddrMode uint

const(
	Immediate AddrMode = iota
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

const(
	NMI = iota
	RESET
	IRQ
)

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
	return c.P & Negative != 0
}

func (c *Cpu) updateNZ(b byte){
	if b == 0x00{
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

func (c *Cpu) updateV(prev byte, now byte){
	if prev & 0x80 == 0 && now & 0x80 != 0{
		c.setBit(Overflow)
	} else{
		c.unsetBit(Overflow)
	}
}

func (c *Cpu) updateC(prev byte, now byte){
	if prev & 0x80 != 0 && now & 0x80 == 0{
		c.setBit(Carry)
	}else{
		c.unsetBit(Carry)
	}
}

func (c *Cpu) lda(b byte){
	c.A = b
	c.updateNZ(c.A)
}

func (c *Cpu) ldx(b byte){
	c.X = b
	c.updateNZ(c.X)
}

func (c *Cpu) ldy(b byte){
	c.Y = b
	c.updateNZ(c.Y)
}

func (c *Cpu) sta(addr Word){
	c.bus.Store(addr, c.A)
}

func (c *Cpu) stx(addr Word){
	c.bus.Store(addr, c.X)
}

func (c *Cpu) sty(addr Word){
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
	c.updateNZ(c.S)
}

func (c *Cpu) tya(){
	c.A = c.Y
	c.updateNZ(c.A)
}

func (c *Cpu) adc(b byte){
	prev := c.A
	c.A = c.A + b + c.status(Carry)
	c.updateC(prev, c.A)
	c.updateV(prev, c.A)
	c.updateNZ(c.A)
}

func (c *Cpu) and(b byte){
	c.A = c.A & b
	c.updateNZ(c.A)
}

func (c *Cpu) asl(isAccumulator bool, addr Word){
	if isAccumulator{
		prev := c.A
		c.A <<= 1
		c.updateC(prev, c.A)
		c.updateNZ(c.A)
	}else{
		v := c.bus.Load(addr)
		c.bus.Store(addr, v << 1)
		c.updateC(v, v << 1)
		c.updateNZ(v)
	}
}

func (c *Cpu) bit(addr Word){
	// 特殊なレジスタ操作が必要なのでロジックを個別化する
	v := c.bus.Load(addr)
	if v >> 6 & 0x01 != 0{
		c.setBit(Overflow)
	}else{
		c.unsetBit(Overflow)
	}

	if v & c.A == 0{
		c.setBit(Zero)
	}else{
		c.unsetBit(Zero)
	}

	if v & 0x80 != 0{
		c.setBit(Negative)
	} else{
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

func (c *Cpu) cmp(addr Word){
	v := c.bus.Load(addr)
	c.compare(c.A, v)
}

func (c *Cpu) cpx(addr Word){
	v := c.bus.Load(addr)
	c.compare(c.X, v)
}

func (c *Cpu) cpy(addr Word){
	v := c.bus.Load(addr)
	c.compare(c.Y, v)
}

func (c *Cpu) dec(addr Word){
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

func (c *Cpu) eor(b byte){
	c.A ^= b
	c.updateNZ(c.A)
}

func (c *Cpu) inc(addr Word){
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

func (c *Cpu) lsr(isAccumulator bool, addr Word){
	if isAccumulator{
		a:=c.A
		c.A<<=1
		c.updateC(a, c.A)
		c.updateNZ(c.A)
	}else{
		v := c.bus.Load(addr)
		c.bus.Store(addr, v>>1)
		c.updateC(v, v>>1)
		c.updateNZ(v)
	}
}

func (c *Cpu) ora(b byte){
	c.A = c.A | b
	c.updateNZ(c.A)
}

func (c *Cpu) rol(isAccumulator bool, addr Word){
	rotateLeft := func(b byte) byte{
		return (b << 1 & 0xFE) | (b >> 7)
	}
	if isAccumulator{
		prev := c.A
		c.A = rotateLeft(c.A)
		c.updateC(prev, c.A)
		c.updateNZ(c.A)
	}else{
		prev := c.bus.Load(addr)
		v := rotateLeft(prev)
		c.bus.Store(addr, v)
		c.updateC(prev, v)
		c.updateNZ(v)
	}
}

func (c *Cpu) ror(isAccumulator bool, addr Word) {
	rotateRight := func(b byte) byte {
		return b >> 1 | (b << 7 & 0x80)
	}
	if isAccumulator {
		prev := c.A
		c.A = rotateRight(c.A)
		c.updateC(prev, c.A)
		c.updateNZ(c.A)
	} else {
		prev := c.bus.Load(addr)
		v := rotateRight(prev)
		c.bus.Store(addr, v)
		c.updateC(prev, v)
	}
}

func (c *Cpu) sbc(addr Word){
	prev := c.A
	c.A = c.A - c.bus.Load(addr) - (1 - c.status(Carry))
	c.updateC(prev, c.A)
	c.updateV(prev, c.A)
	c.updateNZ(c.A)
}

func (c *Cpu) push(b byte){
	c.bus.Store(0x100 | Word(c.S), c.A)
	c.S--
}

func (c *Cpu) pushWord(w Word){
	h := byte(w >> 8)
	l := byte(w & 0xFF)
	c.push(h)
	c.push(l)
}

func (c *Cpu) pop() byte{
	c.S++
	return c.bus.Load(0x100 | Word(c.S))
}

func (c *Cpu) popWord() Word {
	h := c.pop()
	l := c.pop()
	return Word(h << 8 | l)
}

func (c *Cpu) pha(){
	c.push(c.A)
}

func (c *Cpu) php(){
	c.push(c.P)
}

func (c *Cpu) pla(){
	c.A = c.pop()
	c.updateNZ(c.A)
}

func (c *Cpu) plp(){
	c.P = c.pop()
}

func (c *Cpu) jmp(addr Word){
	c.PC = addr
}

func (c *Cpu) jsr(addr Word){
	// NEED - 1?
	c.pushWord(c.PC)
	c.jmp(addr)
}

func (c *Cpu) rst(){
	// NEED + 1?
	c.PC = c.popWord()
}

func (c *Cpu) rti() {
	c.P = c.pop()&0xEF | 0x20
	c.PC = c.popWord()
}

func (c *Cpu) branch(offset byte){
	c.PC = Word(int(c.PC) + int(int8(offset)))
}

func (c *Cpu) bcc(offset byte){
	if !c.isCarry(){
		c.branch(offset)
	}
}

func (c *Cpu) bcs(offset byte){
	if c.isCarry(){
		c.branch(offset)
	}
}

func (c *Cpu) beq(offset byte){
	if c.isZero(){
		c.branch(offset)
	}
}

func (c *Cpu) bne(offset byte){
	if !c.isZero(){
		c.branch(offset)
	}
}

func (c *Cpu) bmi(offset byte){
	if c.isNegative(){
		c.branch(offset)
	}
}

func (c *Cpu) bpl(offset byte){
	if !c.isNegative(){
		c.branch(offset)
	}
}

func (c *Cpu) bvc(offset byte){
	if !c.isOverflow(){
		c.branch(offset)
	}
}

func (c *Cpu) bvs(offset byte){
	if !c.isOverflow(){
		c.branch(offset)
	}
}

func (c *Cpu) clc(){
	c.unsetBit(Carry)
}

func (c *Cpu) cld(){
	// Not Implemented
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

func (c *Cpu) brk(){
	c.pushWord(c.PC)
	c.php()
	c.sei()
	c.PC = Irq
}

func (c *Cpu) nop(){
	// nop
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
	case Implied:
		c.PC += 1
	case Immediate, Zeropage, ZeropageX, ZeropageY, Relative:
		c.PC += 2
	case Absolute, AbsoluteX, AbsoluteY, Indirect, IndirectX, IndirectY:
		c.PC += 3
	default:
		abort("panic: unknown addrMode `%s` was called when advance", mode)
	}
}

func (c *Cpu) solveAddrMode(mode AddrMode) Word {
	switch mode {
	case Implied:
		return 0x00
	case Immediate:
		return Word(c.bus.Load(c.PC + 1))
	case Relative:
		return Word(c.bus.Load(c.PC + 1))
	case Absolute:
		return c.bus.Loadw(c.PC + 1)
	case AbsoluteX:
		return c.bus.Loadw(Word(int(c.bus.Loadw(c.PC + 1)) + int(int8(c.X))))
	case AbsoluteY:
		return c.bus.Loadw(Word(int(c.bus.Loadw(c.PC + 1)) + int(int8(c.Y))))
	default:
		abort("panic: unknown addrMode `%s` was called when solving", mode)
	}
	panic("Unable to reach here")
}

func (c *Cpu) execute(inst Instruction, w Word){
	switch inst.mnemonic {
	case "LDA":
		c.lda(byte(w))
	case "LDX":
		c.ldx(byte(w))
	case "LDY":
		c.ldy(byte(w))
	case "STA":
		c.sta(w)
	case "SEI":
		c.sei()
	case "TXS":
		c.txs()
	case "INX":
		c.inx()
	case "DEY":
		c.dey()
	case "BNE":
		c.bne(byte(w))
	case "JMP":
		c.jmp(w)
	case "BRK":
		c.brk()
	default:
		abort("panic: unknown mnemonic `%s` was invoked.", inst.mnemonic)
	}
}

func (c *Cpu) dump(b byte, arg Word, mne string, mode AddrMode){
	fmt.Printf("[PC:0x%4x, A:0x%2x, X:0x%2x, Y:0x%2x] " +
		"%2x, %4x ## %s (%s) \n",
		c.PC, c.A, c.X, c.Y, b, arg, mne, mode)
}
