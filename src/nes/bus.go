package nes

type Bus struct{
	wram Mem
	cpu *Cpu
	ppu *Ppu
	controller *Controller
	prgRom []byte
}

func NewBus(wram Mem, prgRom []byte) *Bus{
	return &Bus{
		wram: wram,
		prgRom:prgRom,
	}
}

func (b *Bus) Load(addr word) byte{
	if addr < 0x0800 {
		return b.wram.load(addr)
	} else if addr < 0x2000{
	} else if addr == 0x2002 {
		return b.ppu.readPpuStatus()
	} else if addr == 0x2004 {
		return b.ppu.readOamData()
	} else if addr == 0x2007{
		return b.ppu.readPpuData()
	} else if addr == 0x4016{
		// 1P
		return b.controller.read()
	} else if addr >= 0x8000 {
		return b.prgRom[addr - 0x8000]
	}

	abort("[Load] Not implementd address 0x%x", addr)
	return 0x00
}

func (b *Bus) Loadw(addr word) word {
	// little endian
	upper := word(b.Load(addr+1))
	bottom := word(b.Load(addr))
	return upper << 8 | bottom
}

func (b *Bus) Store(addr word, v byte){
	if addr < 0x0800 {
		b.wram.store(addr, v)
	} else if addr < 0x2000 {
		b.wram.store(addr-0x0800, v)
	} else if addr == 0x2000 {
		b.ppu.writePpuCtrl(v)
	} else if addr == 0x2001 {
		b.ppu.writePpuMask(v)
	} else if addr == 0x2003 {
		b.ppu.writeOamAddr(v)
	} else if addr == 0x2004 {
		b.ppu.writeOamData(v)
	} else if addr == 0x2005 {
		// TODO
	} else if addr == 0x2006 {
		b.ppu.writePpuAddr(v)
	} else if addr == 0x2007 {
		b.ppu.writePpuData(v)
	} else if addr == 0x4016{
		// 1P
		b.controller.reset(v)
	} else{
		abort("[Store] Not implementd address 0x%x", addr)
	}
}
