package nes

import (
	"fmt"
	"io/ioutil"
)

type Ines interface {
	PrgRom() []byte
	ChrRom() []byte
}

const HeaderSize = 0x0010

func NewCassette(path string) (Ines, error){
	bytes, err := ioutil.ReadFile(path)
	if err != nil{
		return nil, fmt.Errorf("cannot open the cassete. [PATH] %s", path)
	}

	// validate header
	if !(bytes[0] == 0x4e && bytes[1] == 0x45 && bytes[2] == 0x53 && bytes[3] == 0x1A){
		return nil, fmt.Errorf("validation magic error :this file is not nes file. [PATH] %s", path)
	}

	// prgRom = 0x4000 Byte (16KB) * header[4]
	// chrRom = 0x2000 Byte (8KB) * header[5]
	prgRomStart := HeaderSize
	chrRomStart := HeaderSize + int(bytes[4]) * 0x4000
	chrROMEnd := chrRomStart + int(bytes[5]) * 0x2000

	return &Mapper0{
		prgRom:bytes[prgRomStart:chrRomStart-1],
		chrRom:bytes[chrRomStart:chrROMEnd-1],
	}, nil
}

/// Implement of Mapper0

type Mapper0 struct{
	prgRom []byte
	chrRom []byte
}

func (m *Mapper0) PrgRom() []byte{
	return m.prgRom
}

func (m *Mapper0) ChrRom() []byte{
	return m.chrRom
}
