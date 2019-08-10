package nes

type Ram struct{
	data []byte
}

func NewRam(size int) Mem {
	return &Ram{
		data:make([]byte, size),
	}
}

func NewRamInit(size int, init []byte) Mem {
	ram := make([]byte, size)
	copy(ram, init)
	return &Ram{
		data:ram,
	}
}

func (m *Ram) load(addr Word) byte{
	return m.data[addr]
}

func (m *Ram) store(addr Word, b byte){
	m.data[addr] = b
}

func (m *Ram) slice(begin int, end int) []byte{
	return m.data[begin:end]
}