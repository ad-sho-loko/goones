package nes

type Bus struct{
	wram Mem
	cpu *Cpu
	ppu *Ppu
	prgRom []byte
}

func NewBus(wram Mem, prgRom []byte) *Bus{
	return &Bus{
		wram: wram,
		prgRom:prgRom,
	}
}

func (b *Bus) Load(addr Word) byte{
	if addr < 0x0800 {
		return b.wram.load(addr)
	} else if addr < 0x2000{
		return b.wram.load(addr - 0x8000)
	} else if addr == 0x2007{
		return b.ppu.readPpuData()
	} else if addr >= 0x8000{
		return b.prgRom[addr - 0x8000]
	}
	return 0x00
}

func (b *Bus) Loadw(addr Word) Word {
	// little endian
	upper := Word(b.Load(addr+1))
	bottom := Word(b.Load(addr))
	return upper << 8 | bottom
}

func (b *Bus) Store(addr Word, v byte){
	if addr < 0x0800 {
		b.wram.store(addr, v)
	} else if addr < 0x2000 {
		b.wram.store(addr-0x0800, v)
	} else if addr == 0x2006 {
		b.ppu.writePpuAddr(v)
	} else if addr == 0x2007 {
		b.ppu.writePpuData(v)
	}
}
