package nes

// RAM for cpu
type Ram struct{
	data []byte
}

func NewRam(size int) Mem {
	return &Ram{
		data:make([]byte, size),
	}
}

func (m *Ram) load(addr word) byte{
	return m.data[addr]
}

func (m *Ram) store(addr word, b byte){
	m.data[addr] = b
}

func (m *Ram) slice(begin int, end int) []byte{
	return m.data[begin:end]
}

// Ram for ppu
type VRam struct{
	data []byte
	isHorizontalMirror bool
}

func NewVRamInit(size int, init []byte, isHorizontalMirror bool) Mem {
	ram := make([]byte, size)
	copy(ram, init)
	return &VRam{
		data:ram,
		isHorizontalMirror:isHorizontalMirror,
	}
}

func isNameTable2(addr word) bool{
	return addr >= 0x2400 && addr < 0x2800
}

func isNameTable4(addr word) bool{
	return addr >= 0x2C00 && addr < 0x3000
}

func (m *VRam) load(addr word) byte{
	// always mirroring
	if addr >= 0x3000 && addr < 0x3F00 {
		return m.data[addr - 0x1000]
	}

	if addr == 0x3F10 || addr == 0x3F14 || addr == 0x3F18 || addr == 0x3F1C{
		return m.data[addr - 0x10]
	}

	if addr >= 0x3F20 && addr <= 0x3FFF {
		return m.data[addr - (addr % 0x20)]
	}

	if addr >= 0x4000{
		return m.data[addr%0x4000]
	}

	if isNameTable2(addr) || isNameTable4(addr){
		if m.isHorizontalMirror{
			return m.data[addr-0x0400]
		}else{
			return m.data[addr-0x0800]
		}
	}

	return m.data[addr]
}

func (m *VRam) store(addr word, b byte){
	// always mirroring
	if addr >= 0x3000 && addr < 0x3F00 {
		// 0x3000 - 0x3EFF is mirror of 0x2000 - 0x2EFF
		m.data[addr - 0x1000] = b
		return
	}

	if addr == 0x3F10 || addr == 0x3F14 || addr == 0x3F18 || addr == 0x3F1C{
		// $3F10/$3F14/$3F18/$3F1C are mirror of $3F00/$3F04/$3F08/$3F0C.
		m.data[addr - 0x10] = b
		return
	}

	if addr >= 0x3F20 && addr <= 0x3FFF {
		// 0x3F20 - 0x3FFF is mirror of 0x3F00 - 0x3F1F
		m.data[addr - (addr % 0x20)] = b
		return
	}

	if addr >= 0x4000{
		m.data[addr % 0x4000] = b
		return
	}


	// 水平ミラーリング or 垂直ミラーリング
	if isNameTable2(addr) || isNameTable4(addr){
		if m.isHorizontalMirror{
			m.data[addr-0x0400] = b
		}else{
			m.data[addr-0x0800] = b
		}
		return
	}

	m.data[addr] = b
}

func (m *VRam) slice(begin int, end int) []byte{
	return m.data[begin:end]
}

