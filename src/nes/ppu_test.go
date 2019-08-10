package nes

import "testing"

var p *Ppu

func Test_WritePpuAddr(t *testing.T){
	p.writePpuAddr(0xFF)
	if p.PpuAddr != 0x00FF{
		t.Fatal()
	}

	p.writePpuAddr(0xEE)
	if p.PpuAddr != 0xFFEE{
		t.Fatal()
	}

	p.writePpuAddr(0xDD)
	if p.PpuAddr != 0xEEDD{
		t.Fatal()
	}
}

